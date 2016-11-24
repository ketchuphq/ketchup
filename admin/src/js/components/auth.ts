interface User {
  email: string;
  uuid: string;
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
