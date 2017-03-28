import * as chai from 'chai';
import * as chaiAsPromised from 'chai-as-promised';

// suppress errors
declare var global: any;
declare var Object: any;
declare var require: any;

global.window = Object.assign(
  require('mithril/test-utils/domMock.js')(),
  require('mithril/test-utils/pushStateMock')()
)

// uncomment for extra debugging
// chai.config.includeStack = true;
chai.use(chaiAsPromised);
chai.use(require('chai-subset'));