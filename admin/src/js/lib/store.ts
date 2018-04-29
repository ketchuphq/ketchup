export interface Data<T> {
  initial?: T;
  current?: T;
}

export default class GenericStore<T> {
  obj: T;
  copy: (from: T, to?: T) => T;
  private subscribers: {[key: string]: (obj: T) => void};

  constructor(copy: (from: T, to?: T) => T, obj?: T) {
    this.obj = obj;
    this.copy = copy;
    this.subscribers = {};
  }

  subscribe(key: string, fn: (obj: T) => void) {
    this.subscribers[key] = fn;
  }

  unsubscribe(key: string) {
    delete this.subscribers[key];
  }

  set(obj: T): T {
    this.obj = this.copy(obj);
    this.notify();
    return this.obj;
  }

  update(fn: (v: T) => void): T {
    fn(this.obj);
    return this.set(this.obj);
  }

  protected notify() {
    Object.keys(this.subscribers).map((key) => {
      this.subscribers[key](this.obj);
    });
  }
}
