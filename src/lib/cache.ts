// Cache Policy Simulator - Pure TypeScript Implementation

export interface CacheInterface {
  get(key: number): number | undefined;
  put(key: number, value: number): number | null;
  size(): number;
  capacity(): number;
  clear(): void;
  keys(): number[];
  getName(): string;
}

// FIFO (First In First Out) Cache
export class FIFOCache implements CacheInterface {
  private cap: number;
  private elements: Map<number, number>;
  private order: number[];

  constructor(capacity: number) {
    this.cap = Math.max(1, capacity);
    this.elements = new Map();
    this.order = [];
  }

  get(key: number): number | undefined {
    return this.elements.get(key);
  }

  put(key: number, value: number): number | null {
    if (this.elements.has(key)) {
      this.elements.set(key, value);
      return null;
    }

    let evicted: number | null = null;

    if (this.order.length >= this.cap) {
      evicted = this.order.shift()!;
      this.elements.delete(evicted);
    }

    this.elements.set(key, value);
    this.order.push(key);

    return evicted;
  }

  size(): number {
    return this.elements.size;
  }

  capacity(): number {
    return this.cap;
  }

  clear(): void {
    this.elements.clear();
    this.order = [];
  }

  keys(): number[] {
    return [...this.order];
  }

  getName(): string {
    return 'FIFO';
  }
}

// LRU (Least Recently Used) Cache
export class LRUCache implements CacheInterface {
  private cap: number;
  private elements: Map<number, number>;
  private order: number[]; // Front = LRU, Back = MRU

  constructor(capacity: number) {
    this.cap = Math.max(1, capacity);
    this.elements = new Map();
    this.order = [];
  }

  get(key: number): number | undefined {
    if (!this.elements.has(key)) {
      return undefined;
    }
    // Move to back (MRU)
    const idx = this.order.indexOf(key);
    if (idx > -1) {
      this.order.splice(idx, 1);
      this.order.push(key);
    }
    return this.elements.get(key);
  }

  put(key: number, value: number): number | null {
    if (this.elements.has(key)) {
      this.elements.set(key, value);
      // Move to back
      const idx = this.order.indexOf(key);
      if (idx > -1) {
        this.order.splice(idx, 1);
        this.order.push(key);
      }
      return null;
    }

    let evicted: number | null = null;

    if (this.order.length >= this.cap) {
      evicted = this.order.shift()!;
      this.elements.delete(evicted);
    }

    this.elements.set(key, value);
    this.order.push(key);

    return evicted;
  }

  size(): number {
    return this.elements.size;
  }

  capacity(): number {
    return this.cap;
  }

  clear(): void {
    this.elements.clear();
    this.order = [];
  }

  keys(): number[] {
    return [...this.order];
  }

  getName(): string {
    return 'LRU';
  }
}

// LFU (Least Frequently Used) Cache
interface LFUEntry {
  key: number;
  value: number;
  freq: number;
}

export class LFUCache implements CacheInterface {
  private cap: number;
  private elements: Map<number, LFUEntry>;
  private freqLists: Map<number, number[]>; // freq -> list of keys (ordered by insertion)
  private minFreq: number;

  constructor(capacity: number) {
    this.cap = Math.max(1, capacity);
    this.elements = new Map();
    this.freqLists = new Map();
    this.minFreq = 0;
  }

  private incrementFreq(key: number): void {
    const entry = this.elements.get(key)!;
    const oldFreq = entry.freq;
    const newFreq = oldFreq + 1;

    // Remove from old freq list
    const oldList = this.freqLists.get(oldFreq);
    if (oldList) {
      const idx = oldList.indexOf(key);
      if (idx > -1) {
        oldList.splice(idx, 1);
      }
      if (oldList.length === 0) {
        this.freqLists.delete(oldFreq);
        if (this.minFreq === oldFreq) {
          this.minFreq = newFreq;
        }
      }
    }

    // Add to new freq list
    if (!this.freqLists.has(newFreq)) {
      this.freqLists.set(newFreq, []);
    }
    this.freqLists.get(newFreq)!.push(key);
    entry.freq = newFreq;
  }

  get(key: number): number | undefined {
    const entry = this.elements.get(key);
    if (!entry) {
      return undefined;
    }
    this.incrementFreq(key);
    return entry.value;
  }

  put(key: number, value: number): number | null {
    if (this.elements.has(key)) {
      const entry = this.elements.get(key)!;
      entry.value = value;
      this.incrementFreq(key);
      return null;
    }

    let evicted: number | null = null;

    if (this.elements.size >= this.cap) {
      // Evict from min freq list (oldest first)
      const minList = this.freqLists.get(this.minFreq);
      if (minList && minList.length > 0) {
        evicted = minList.shift()!;
        this.elements.delete(evicted);
        if (minList.length === 0) {
          this.freqLists.delete(this.minFreq);
        }
      }
    }

    // Add new entry with freq 1
    this.elements.set(key, { key, value, freq: 1 });
    if (!this.freqLists.has(1)) {
      this.freqLists.set(1, []);
    }
    this.freqLists.get(1)!.push(key);
    this.minFreq = 1;

    return evicted;
  }

  size(): number {
    return this.elements.size;
  }

  capacity(): number {
    return this.cap;
  }

  clear(): void {
    this.elements.clear();
    this.freqLists.clear();
    this.minFreq = 0;
  }

  keys(): number[] {
    return Array.from(this.elements.keys());
  }

  getName(): string {
    return 'LFU';
  }
}

// Random Replacement Cache
export class RandomCache implements CacheInterface {
  private cap: number;
  private elements: Map<number, number>;
  private keyList: number[];

  constructor(capacity: number) {
    this.cap = Math.max(1, capacity);
    this.elements = new Map();
    this.keyList = [];
  }

  get(key: number): number | undefined {
    return this.elements.get(key);
  }

  put(key: number, value: number): number | null {
    if (this.elements.has(key)) {
      this.elements.set(key, value);
      return null;
    }

    let evicted: number | null = null;

    if (this.keyList.length >= this.cap) {
      const randomIndex = Math.floor(Math.random() * this.keyList.length);
      evicted = this.keyList[randomIndex];
      this.elements.delete(evicted);
      // Swap with last and pop for O(1) removal
      this.keyList[randomIndex] = this.keyList[this.keyList.length - 1];
      this.keyList.pop();
    }

    this.elements.set(key, value);
    this.keyList.push(key);

    return evicted;
  }

  size(): number {
    return this.elements.size;
  }

  capacity(): number {
    return this.cap;
  }

  clear(): void {
    this.elements.clear();
    this.keyList = [];
  }

  keys(): number[] {
    return [...this.keyList];
  }

  getName(): string {
    return 'Random';
  }
}

// Factory function
export function createCache(policy: string, capacity: number): CacheInterface {
  switch (policy) {
    case 'FIFO':
      return new FIFOCache(capacity);
    case 'LRU':
      return new LRUCache(capacity);
    case 'LFU':
      return new LFUCache(capacity);
    case 'Random':
      return new RandomCache(capacity);
    default:
      return new FIFOCache(capacity);
  }
}
