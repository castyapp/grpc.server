package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

var (
	ErrKeyNodFound = errors.New("key not found")
)

type Provider interface {
	Register(ctx *Context) error
	Close(ctx *Context) error
}

type Context struct {
	ctx          context.Context
	items        map[string]interface{}
	providers    []Provider
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

func NewContext(ctx context.Context) *Context {
	return &Context{ctx: ctx, items: make(map[string]interface{})}
}

func (c *Context) With(handlers ...Provider) *Context {
	for _, handler := range handlers {
		if err := handler.Register(c); err != nil {
			fmt.Printf("Could not register provider [%s] cause [%v]\n", reflect.ValueOf(handler).Kind().String(), err)
			os.Exit(1)
		}
		c.providers = append(c.providers, handler)
	}
	return c
}

func (c *Context) Close() error {
	for _, p := range c.providers {
		if err := p.Close(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Context) MustGet(key string) (value interface{}) {
	value, err := c.Get(key)
	if err != nil {
		log.Panicln(err)
	}
	return value
}

func (c *Context) MustGetBoolean(key string) bool {
	value, err := c.Get(key)
	if err != nil {
		log.Panicln(err)
	}
	return value.(bool)
}

func (c *Context) GetBoolean(key string) (bool, error) {
	value, err := c.Get(key)
	if err != nil {
		return false, err
	}
	return value.(bool), nil
}

func (c *Context) GetString(key string) (string, error) {
	value, err := c.Get(key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (c *Context) MustGetString(key string) string {
	value, err := c.Get(key)
	if err != nil {
		log.Panicln(err)
	}
	return value.(string)
}

func (c *Context) Count() int {
	count := 0
	c.RLock()
	count += len(c.items)
	c.RUnlock()
	return count
}

func (c *Context) Get(key string) (interface{}, error) {
	c.RLock()
	val, ok := c.items[key]
	if !ok {
		return nil, ErrKeyNodFound
	}
	c.RUnlock()
	return val, nil
}

func (c *Context) Set(key string, value interface{}) error {
	c.Lock()
	if _, ok := c.items[key]; ok {
		return fmt.Errorf(fmt.Sprintf("Key [%s] already exists!", key))
	}
	c.items[key] = value
	c.Unlock()
	return nil
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}
