import * as React from 'react';
import * as ReactDOM from 'react-dom';
import GenericStore from 'lib/store';

const sweepInterval = 500;
let ToastStore = new GenericStore<Toast[]>((from, _?) => from.slice(), []);

type ToastType = 'error' | 'green';

export default class Toaster extends React.Component {
  render() {
    return ReactDOM.createPortal(<ToastList />, document.getElementById('react-root'));
  }
}

class ToastList extends React.Component<{}, {toasts: Toast[]}> {
  componentWillMount() {
    ToastStore.subscribe('ToastList', (toasts) => {
      this.setState({toasts});
    });
  }

  componentWillUnmount() {
    ToastStore.unsubscribe('ToastList');
  }

  render() {
    let toasts = ToastStore.obj;
    let klass = 'toast-wrapper';
    if (toasts.length == 0) {
      klass += ' toast-wrapper--hidden';
    }
    return (
      <div className={klass}>
        {toasts.map((t) => {
          let k = ['toast', 'toast--enter', t.getClass(), t.expired() ? 'toast--expired' : ''];
          return (
            <div key={t.key} className={k.join(' ')}>
              <div className="contents">{t.message}</div>
            </div>
          );
        })}
      </div>
    );
  }
}

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
  let toasts = ToastStore.obj;
  let now = new Date().getTime() - 2000;
  let prev = toasts.length;
  let hasChange = toasts.reduce((prev, cur) => {
    return prev || cur.expired() != cur.expired(now);
  }, false);
  toasts = toasts.filter((t) => !t.expired(now));
  if (prev != toasts.length || hasChange) {
    ToastStore.set(toasts);
  }
}, sweepInterval);

export function add(message: string, klass: ToastType = 'green') {
  ToastStore.update((toasts) => {
    toasts.push(new Toast(message, 3000, klass));
  });
}

export const error = (message: string) => add(message, 'error');
