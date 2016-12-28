import * as store from 'store';
import * as Toaster from 'components/toaster';

export interface User {
  email: string;
  uuid: string;
}

interface Preferences {
  hideMenu: boolean;
}

let cachedUser: User = null;

// AuthController is a super class for controllers which may require auth
export class AuthController {
  user: Mithril.BasicProperty<User>;
  _userPromise: Mithril.Promise<User>;

  constructor(user: User = null) {
    this.user = m.prop<User>(user || cachedUser);
    if (this.user()) {
      var deferred = m.deferred<User>();
      deferred.resolve(this.user());
      this._userPromise = deferred.promise;
    } else {
      this._userPromise = m.request({
        method: 'GET',
        url: '/api/v1/user',
        background: false,
      }).then((res: User) => {
        if (!res.uuid) {
          throw new Error('not logged in');
        }
        this.user(res);
        cachedUser = res;
        return res;
      }).catch(() => {
        cachedUser = null;
        return null;
      });
    }
  }

  logout() {
    m.request({
      method: 'GET',
      url: '/api/v1/logout',
      background: false
    }).then(() => {
      Toaster.add('logged out');
      setTimeout(() => {
        location.reload();
      }, 2000);
    });
  }

  private get storeKey(): string {
    if (!this.user()) {
      return null;
    }
    return `user-${this.user().uuid}`;
  }

  private getPrefs(): Preferences {
    if (!this.user()) {
      return null;
    }
    let prefs: Preferences = store.get(this.storeKey);
    if (!prefs) {
      prefs = {
        hideMenu: false,
      }
      store.set(this.storeKey, prefs);
    }
    return prefs;
  }

  pref<K extends keyof Preferences>(key: K): Preferences[K] {
    if (!this.user()) {
      return; // error
    }
    let prefs: Preferences = this.getPrefs();
    if (!prefs) {
      return; // error
    }
    return prefs[key];
  }

  setPref<K extends keyof Preferences>(key: K, val: Preferences[K]) {
    if (!this.user()) {
      return; // error
    }
    let prefs: Preferences = this.getPrefs();
    prefs[key] = val;
    store.set(this.storeKey, prefs);
  }
}

export class MustAuthController extends AuthController {
  ready: boolean;
  constructor() {
    super();
    this.ready = false;
    this._userPromise.then((user) => {
      if (!user || !user.uuid) {
        window.location.assign(`/admin/login?next=${window.location.pathname}`);
      } else {
        this.ready = true;
      }
    });
  }
}
