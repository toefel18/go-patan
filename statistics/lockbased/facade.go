/*
 *
 *     Copyright 2016 Christophe Hesters
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 *
 */
package lockbased

import (
    "github.com/toefel18/go-patan/statistics/api"
)

type Facade struct {
    store *Store
}

func NewFacade(store *Store) *Facade {
    if store == nil {
        panic("store = nil, Facade needs a store")
    }
    return &Facade{store}
}

func (facade *Facade) StartStopwatch() api.Stopwatch {
    return startNewStopwatch()
}

func (facade *Facade) RecordElapsedTime(key string, stopwatch api.Stopwatch) int64 {
    millis := stopwatch.ElapsedMillis()
    facade.store.addDuration(key, millis)
    return millis
}

func (facade *Facade) MeasureFunc(key string, subject func()) int64 {
    sw := startNewStopwatch()
    subject()
    return facade.RecordElapsedTime(key, sw)
}

func (facade *Facade) IncrementCounter(key string) {
    facade.AddToCounter(key, 1)
}

func (facade *Facade) DecrementCounter(key string) {
    facade.AddToCounter(key, -1)
}

func (facade *Facade) AddToCounter(key string, value int64) {
    facade.store.addToCounter(key, value)
}

func (facade *Facade) AddSample(key string, value int64) {
    facade.store.addSample(key, value)
}

func (facade *Facade) Reset() {
    facade.store.Reset();
}

func (facade *Facade) Snapshot() api.Snapshot {
    return facade.store.Snapshot()
}

func (facade *Facade) SnapshotAndReset() api.Snapshot {
    return facade.store.SnapshotAndReset()
}

func (facade *Facade) Close() {
    // nop
}