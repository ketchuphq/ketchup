import msx from 'lib/msx';
import * as m from 'mithril';
import Theme from 'lib/theme';
import { MustAuthController, BaseComponent } from 'components/auth';
import Button from 'components/button';
import { ModalAttrs, ModalComponent } from 'components/modal';
import { Package } from 'lib/api';
import { Table, Row, LinkRow } from 'components/table';
import { ConfirmModalComponent } from 'components/modal';

interface ThemeProps {
  theme: Theme;
}

class VersionSection extends BaseComponent<ThemeProps> {
  showUpdateModal: boolean;
  hasUpdate: boolean;
  latestRef?: string;

  constructor(v: any) {
    super(v);
    this.showUpdateModal = false;
  }

  checkUpdates() {
    return this.props.theme.checkForUpdates().then(({ currentRef }) => {
      if (!currentRef) {
        this.latestRef = 'No updates found.';
        this.hasUpdate = false;
      } else {
        this.latestRef = `Latest ref: ${currentRef.slice(0, 6)}`;
        this.hasUpdate = true;
      }
      this.showUpdateModal = true;
      m.redraw();
    });
  }

  updateTheme() {
    return Promise.resolve();
  }

  view() {
    let checkUpdates;
    if (!this.props.theme.ref) {
      return;
    }
    return [
      <h2>Version</h2>,
      <div class='table'>
        <div class='tr tr--center'>
          <code>{this.props.theme.ref}</code>
          <Button class='button--green button--small' handler={() => this.checkUpdates()}>
            Check for Updates
          </Button>
        </div>
      </div>,
      <ConfirmModalComponent
        title='Updates'
        visible={() => this.showUpdateModal}
        toggle={() => {
          this.showUpdateModal = !this.showUpdateModal;
          m.redraw();
        }}
        confirmText={this.hasUpdate ? 'Update' : 'Okay'}
        hideCancel={this.hasUpdate}
        resolve={this.updateTheme.bind(this)}
      >
        <p>{this.latestRef}</p>
      </ConfirmModalComponent>
    ];
  }
}

class TemplatesSection extends BaseComponent<ThemeProps> {
  view() {
    let theme = this.props.theme;
    let templateKeys = Object.keys(theme.templates);
    let templates = templateKeys.sort().map((name) => theme.templates[name]).map((t) => (
      <LinkRow href={`/admin/themes/${theme.name}/templates/${t.name}`} link>
        <div>{t.name}</div>
        <div>{t.engine}</div>
      </LinkRow>
    ));

    return [<h2>Templates</h2>, <div class='table'>{templates}</div>];
  }
}

class AssetsSection extends BaseComponent<ThemeProps> {
  view() {
    let assetKeys = Object.keys(this.props.theme.assets);
    let assetsList: m.Children = <a class='tr'>no assets</a>;
    if (assetKeys.length > 0) {
      assetsList = assetKeys.sort().map((asset) => <LinkRow href={`/${asset}`}>{asset}</LinkRow>);
    }
    return [<h2>Assets</h2>, <Table>{assetsList}</Table>];
  }
}

interface PackageProps {
  pkg: Package;
}

class PackageSection extends BaseComponent<PackageProps> {
  view() {
    if (!this.props.pkg) {
      return;
    }
    let pkg = this.props.pkg;
    let fields = [];
    if (pkg.authors && pkg.authors.length > 0) {
      fields.push(
        <Row>
          <div>Authors</div>
          <div>{pkg.authors.join(', ')}</div>
        </Row>
      );
    }

    if (pkg.homepage) {
      fields.push(
        <Row>
          <div>Homepage</div>
          <div>{pkg.homepage}</div>
        </Row>
      );
    }
    if (pkg.vcsUrl) {
      fields.push(
        <Row>
          <div>Source</div>
          <div>{pkg.vcsUrl}</div>
        </Row>
      );
    }
    if (pkg.tags && pkg.tags.length > 0) {
      fields.push(
        <Row>
          <div>Tags</div>
          <div>{pkg.tags.join(', ')}</div>
        </Row>
      );
    }

    if (fields.length == 0) {
      return;
    }

    return [<h2>Package</h2>, <Table>{fields}</Table>];
  }
}

export default class ThemePage extends MustAuthController {
  theme: Theme;

  constructor() {
    super();
    let themeName = m.route.param('name');
    if (themeName) {
      Theme.get(themeName).then((theme) => {
        this.theme = theme;
        m.redraw();
      });
    }
  }

  view() {
    if (!this.theme) {
      return;
    }
    return (
      <div class='theme'>
        <header>
          <h1>
            <a href='/admin/themes' oncreate={m.route.link}>
              Themes
            </a>
            {m.trust(' &rsaquo; ')}
            <span class='unbold'>{this.theme.name}</span>
          </h1>
        </header>
        <p class='txt-large'>{this.theme.description}</p>

        <VersionSection theme={this.theme} />
        <PackageSection pkg={this.theme.package} />
        <TemplatesSection theme={this.theme} />
        <AssetsSection theme={this.theme} />
      </div>
    );
  }
}
