/*
package apisrv provides an implementation of the gRPC server defined in ../../../api/protobuf-spec/frontend.proto.

Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package apisrv

import (
	"context"
	"errors"
	"net"
	"time"

	frontend "github.com/GoogleCloudPlatform/open-match/cmd/frontendapi/proto"
	"github.com/GoogleCloudPlatform/open-match/internal/metrics"
	playerq "github.com/GoogleCloudPlatform/open-match/internal/statestorage/redis/playerq"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

// Logrus structured logging setup
var (
	feLogFields = log.Fields{
		"app":       "openmatch",
		"component": "frontend",
		"caller":    "frontendapi/apisrv/apisrv.go",
	}
	feLog = log.WithFields(feLogFields)
)

// FrontendAPI implements frontend.ApiServer, the server generated by compiling
// the protobuf, by fulfilling the frontend.APIClient interface.
type FrontendAPI struct {
	grpc *grpc.Server
	cfg  *viper.Viper
	pool *redis.Pool
}
type frontendAPI FrontendAPI

// New returns an instantiated srvice
func New(cfg *viper.Viper, pool *redis.Pool) *FrontendAPI {
	s := FrontendAPI{
		pool: pool,
		grpc: grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{})),
		cfg:  cfg,
	}

	// Add a hook to the logger to auto-count log lines for metrics output thru OpenCensus
	log.AddHook(metrics.NewHook(FeLogLines, KeySeverity))

	// Register gRPC server
	frontend.RegisterAPIServer(s.grpc, (*frontendAPI)(&s))
	feLog.Info("Successfully registered gRPC server")
	return &s
}

// Open opens the api grpc service, starting it listening on the configured port.
func (s *FrontendAPI) Open() error {
	ln, err := net.Listen("tcp", ":"+s.cfg.GetString("api.frontend.port"))
	if err != nil {
		feLog.WithFields(log.Fields{
			"error": err.Error(),
			"port":  s.cfg.GetInt("api.frontend.port"),
		}).Error("net.Listen() error")
		return err
	}
	feLog.WithFields(log.Fields{"port": s.cfg.GetInt("api.frontend.port")}).Info("TCP net listener initialized")

	go func() {
		err := s.grpc.Serve(ln)
		if err != nil {
			feLog.WithFields(log.Fields{"error": err.Error()}).Error("gRPC serve() error")
		}
		feLog.Info("serving gRPC endpoints")
	}()

	return nil
}

// CreateRequest is this service's implementation of the CreateRequest gRPC method // defined in ../proto/frontend.proto
func (s *frontendAPI) CreateRequest(c context.Context, g *frontend.Group) (*frontend.Result, error) {

	// Get redis connection from pool
	redisConn := s.pool.Get()
	defer redisConn.Close()

	// Create context for tagging OpenCensus metrics.
	funcName := "CreateRequest"
	fnCtx, _ := tag.New(c, tag.Insert(KeyMethod, funcName))

	// Write group
	// TODO: Remove playerq module and just use redishelper module once
	// indexing has its own implementation
	err := playerq.Create(redisConn, g.Id, g.Properties)

	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Error("State storage error")

		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.Result{Success: false, Error: err.Error()}, err
	}

	stats.Record(fnCtx, FeGrpcRequests.M(1))
	return &frontend.Result{Success: true, Error: ""}, err

}

// DeleteRequest is this service's implementation of the DeleteRequest gRPC method defined in
// frontendapi/proto/frontend.proto
func (s *frontendAPI) DeleteRequest(c context.Context, g *frontend.Group) (*frontend.Result, error) {
	// Get redis connection from pool
	redisConn := s.pool.Get()
	defer redisConn.Close()

	// Create context for tagging OpenCensus metrics.
	funcName := "DeleteRequest"
	fnCtx, _ := tag.New(c, tag.Insert(KeyMethod, funcName))

	// Write group
	err := playerq.Delete(redisConn, g.Id)
	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Error("State storage error")

		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.Result{Success: false, Error: err.Error()}, err
	}

	stats.Record(fnCtx, FeGrpcRequests.M(1))
	return &frontend.Result{Success: true, Error: ""}, err

}

// GetAssignment is this service's implementation of the GetAssignment gRPC method defined in
// frontendapi/proto/frontend.proto
func (s *frontendAPI) GetAssignment(c context.Context, p *frontend.PlayerId) (*frontend.ConnectionInfo, error) {
	// Get cancellable context
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	// Create context for tagging OpenCensus metrics.
	funcName := "GetAssignment"
	fnCtx, _ := tag.New(ctx, tag.Insert(KeyMethod, funcName))

	// get and return connection string
	var connString string
	watchChan := s.watcher(ctx, s.pool, p.Id) // watcher() runs the appropriate Redis commands.

	select {
	case <-time.After(30 * time.Second): // TODO: Make this configurable.
		err := errors.New("did not see matchmaking results in redis before timeout")
		// TODO:Timeout: deal with the fallout
		// When there is a timeout, need to send a stop to the watch channel.
		// cancelling ctx isn't doing it.
		//cancel()
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
			"playerid":  p.Id,
		}).Error("State storage error")

		errTag, _ := tag.NewKey("errtype")
		fnCtx, _ := tag.New(ctx, tag.Insert(errTag, "watch_timeout"))
		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.ConnectionInfo{ConnectionString: ""}, err

	case connString = <-watchChan:
		feLog.Debug(p.Id, "connString:", connString)
	}

	stats.Record(fnCtx, FeGrpcRequests.M(1))
	return &frontend.ConnectionInfo{ConnectionString: connString}, nil
}

// DeleteAssignment is this service's implementation of the DeleteAssignment gRPC method defined in
// frontendapi/proto/frontend.proto
func (s *frontendAPI) DeleteAssignment(c context.Context, p *frontend.PlayerId) (*frontend.Result, error) {

	// Get redis connection from pool
	redisConn := s.pool.Get()
	defer redisConn.Close()

	// Create context for tagging OpenCensus metrics.
	funcName := "DeleteAssignment"
	fnCtx, _ := tag.New(c, tag.Insert(KeyMethod, funcName))

	// Write group
	err := playerq.Delete(redisConn, p.Id)
	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Error("State storage error")

		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.Result{Success: false, Error: err.Error()}, err
	}

	stats.Record(fnCtx, FeGrpcRequests.M(1))
	return &frontend.Result{Success: true, Error: ""}, err

}

//TODO: Everything below this line will be moved to the redis statestorage library
// in an upcoming version.
// ================================================

// watcher makes a channel and returns it immediately.  It also launches an
// asynchronous goroutine that watches a redis key and returns the value of
// the 'connstring' field of that key once it exists on the channel.
//
// The pattern for this function is from 'Go Concurrency Patterns', it is a function
// that wraps a closure goroutine, and returns a channel.
// reference: https://talks.golang.org/2012/concurrency.slide#25
func (s *frontendAPI) watcher(ctx context.Context, pool *redis.Pool, key string) <-chan string {
	// Add the key as a field to all logs for the execution of this function.
	feLog = feLog.WithFields(log.Fields{"key": key})
	feLog.Debug("Watching key in statestorage for changes")

	watchChan := make(chan string)

	go func() {
		// var declaration
		var results string
		var err = errors.New("haven't queried Redis yet")

		// Loop, querying redis until this key has a value
		for err != nil {
			select {
			case <-ctx.Done():
				// Cleanup
				close(watchChan)
				return
			default:
				results, err = s.retrieveConnstring(ctx, pool, key, s.cfg.GetString("jsonkeys.connstring"))
				if err != nil {
					time.Sleep(5 * time.Second) // TODO: exp bo + jitter
				}
			}
		}
		// Return value retreived from Redis asynchonously and tell calling function we're done
		feLog.Debug("Statestorage watched record update detected")
		watchChan <- results
		close(watchChan)
	}()

	return watchChan
}

// retrieveConnstring is a concurrent-safe, context-aware redis HGET of the 'connstring' fieldin the input key
// TODO: This will be moved to the redis statestorage module.
func (s *frontendAPI) retrieveConnstring(ctx context.Context, pool *redis.Pool, key string, field string) (string, error) {

	// Add the key as a field to all logs for the execution of this function.
	feLog = feLog.WithFields(log.Fields{"key": key})

	cmd := "HGET"
	feLog.WithFields(log.Fields{"query": cmd}).Debug("Statestorage operation")

	// Get a connection to redis
	redisConn, err := pool.GetContext(ctx)
	defer redisConn.Close()

	// Encountered an issue getting a connection from the pool.
	if err != nil {
		feLog.WithFields(log.Fields{
			"error": err.Error(),
			"query": cmd}).Error("Statestorage connection error")
		return "", err
	}

	// Run redis query and return
	return redis.String(redisConn.Do("HGET", key, field))
}
