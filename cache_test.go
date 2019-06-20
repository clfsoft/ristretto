/*
 * Copyright 2019 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ristretto

import (
	"testing"

	"github.com/dgraph-io/ristretto/ring"
)

func TestCache(t *testing.T) {
}

////////////////////////////////////////////////////////////////////////////////

func TestLFU(t *testing.T) {
	t.Run("push", func(t *testing.T) {
		p := NewLFU(4)
		p.Add("1")
		p.Push([]ring.Element{"1", "1", "1"})
		if p.data["1"] != 4 {
			t.Fatal("push error")
		}
	})
	t.Run("add", func(t *testing.T) {
		p := NewLFU(4)
		p.Add("1")
		p.Add("2")
		p.Add("3")
		p.Add("4")
		p.Push([]ring.Element{
			"1", "1", "1", "1",
			"2", "2", "2",
			"3",
			"4", "4",
		})
		victim, added := p.Add("5")
		if added && victim != "3" {
			t.Fatal("eviction error")
		}
	})
}

func BenchmarkLFU(b *testing.B) {
	k := "1"
	data := []ring.Element{"1", "1"}
	b.Run("single", func(b *testing.B) {
		p := NewLFU(1000000)
		p.Add(k)
		b.SetBytes(1)
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			p.hit(k)
		}
	})
	b.Run("parallel", func(b *testing.B) {
		p := NewLFU(1000000)
		p.Add(k)
		b.SetBytes(1)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				p.Push(data)
			}
		})
	})
}
