import msx from 'lib/msx';
import * as m from 'mithril';
const sweepInterval = 500;
let toasts: Toast[] = [];

type ToastType = 'error' | 'green';

class Toast {
  public readonly key: number;
  public entered: boolean;
  private readonly expires: number;
  constructor(public message: string, expiresSeconds: number, public klass: ToastType = 'green') {
    let now = new Date().getTime();
    this.expires = now + expiresSeconds;
    this.key = now;
    this.entered = false;
  }

  expired(now = new Date().getTime()) {
    return this.expires < now;
  }

  getClass() {
    return `toast--${this.klass}`;
  }
}

setInterval(() => {
  let now = new Date().getTime() - 2000;
  let prev = toasts.length;
  let hasChange = toasts.reduce((prev, cur) => {
    return prev || cur.expired() != cur.expired(now);
  }, false);
  toasts = toasts.filter((t) => !t.expired(now));
  if (prev != toasts.length || hasChange) {
    m.redraw();
  }
}, sweepInterval);

export function add(message: string, klass: ToastType = 'green') {
  toasts.push(new Toast(message, 3000, klass));
  m.redraw();
}

export function render() {
  let klass = 'toast-wrapper';
  if (toasts.length == 0) {
    klass += ' toast-wrapper--hidden'
  }
  return <div
    class={klass}
    config={(el: HTMLElement, isInitialized: boolean) => {
      // prevents animation on route change
      if (!isInitialized) {
        for (var i = 0; i < el.children.length; i++) {
          var element = el.children[i];
          element.classList.add('toast--noanimate');
        }
      }
    }}
  >
    {toasts.map((t) =>
      <div
        key={t.key}
        class={['toast toast--enter', t.getClass(), t.expired() ? 'toast--expired' : ''].join(' ')}
      >
        <div class='contents'>{t.message}</div>
      </div>
    )}
  </div>;
}
