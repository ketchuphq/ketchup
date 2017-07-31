// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class Page {
  uuid?: string;
  title?: string;
  theme?: string;
  template?: string;
  timestamps?: Timestamp;
  publishedAt?: string;
  contents?: Content[];
  metadata?: { [key: string]: string; };
  tags?: string[];
  authors?: Author[];
  static copy(from: Page, to?: Page): Page {
    to = to || {};
    to.uuid = from.uuid;
    to.title = from.title;
    to.theme = from.theme;
    to.template = from.template;
    if ('timestamps' in from) {
      to.timestamps = Timestamp.copy(from.timestamps || {}, to.timestamps || {});
    }
    to.publishedAt = from.publishedAt;
    to.contents = from.contents;
    to.metadata = from.metadata;
    to.tags = from.tags;
    to.authors = from.authors;
    return to;
  }
}

export abstract class Content {
  uuid?: string;
  key?: string;
  value?: string;
  timestamps?: Timestamp;
  // skipped field: type
  
  // oneof types:
  multiple?: ContentMultiple;
  short?: ContentString;
  text?: ContentText;
  static copy(from: Content, to?: Content): Content {
    to = to || {};
    to.uuid = from.uuid;
    to.key = from.key;
    to.value = from.value;
    if ('timestamps' in from) {
      to.timestamps = Timestamp.copy(from.timestamps || {}, to.timestamps || {});
    }
    if ('multiple' in from) {
      to.multiple = ContentMultiple.copy(from.multiple || {}, to.multiple || {});
    }
    if ('short' in from) {
      to.short = ContentString.copy(from.short || {}, to.short || {});
    }
    if ('text' in from) {
      to.text = ContentText.copy(from.text || {}, to.text || {});
    }
    return to;
  }
}

export abstract class Route {
  uuid?: string;
  path?: string;
  // skipped field: target
  
  // oneof types:
  file?: string;
  pageUuid?: string;
  static copy(from: Route, to?: Route): Route {
    to = to || {};
    to.uuid = from.uuid;
    to.path = from.path;
    to.file = from.file;
    to.pageUuid = from.pageUuid;
    return to;
  }
}

export abstract class Timestamp {
  createdAt?: string;
  updatedAt?: string;
  static copy(from: Timestamp, to?: Timestamp): Timestamp {
    to = to || {};
    to.createdAt = from.createdAt;
    to.updatedAt = from.updatedAt;
    return to;
  }
}

export abstract class Theme {
  uuid?: string;
  name?: string;
  description?: string;
  package?: Package;
  templates?: { [key: string]: ThemeTemplate; };
  assets?: { [key: string]: ThemeAsset; };
  static copy(from: Theme, to?: Theme): Theme {
    to = to || {};
    to.uuid = from.uuid;
    to.name = from.name;
    to.description = from.description;
    if ('package' in from) {
      to.package = Package.copy(from.package || {}, to.package || {});
    }
    to.templates = from.templates;
    to.assets = from.assets;
    return to;
  }
}

export abstract class ThemeTemplate {
  uuid?: string;
  name?: string;
  theme?: string;
  engine?: string;
  hideContent?: boolean;
  description?: string;
  placeholders?: ThemePlaceholder[];
  data?: string;
  static copy(from: ThemeTemplate, to?: ThemeTemplate): ThemeTemplate {
    to = to || {};
    to.uuid = from.uuid;
    to.name = from.name;
    to.theme = from.theme;
    to.engine = from.engine;
    to.hideContent = from.hideContent;
    to.description = from.description;
    to.placeholders = from.placeholders;
    to.data = from.data;
    return to;
  }
}

export abstract class ThemePlaceholder {
  key?: string;
  // skipped field: type
  
  // oneof types:
  multiple?: ContentMultiple;
  short?: ContentString;
  text?: ContentText;
  static copy(from: ThemePlaceholder, to?: ThemePlaceholder): ThemePlaceholder {
    to = to || {};
    to.key = from.key;
    if ('multiple' in from) {
      to.multiple = ContentMultiple.copy(from.multiple || {}, to.multiple || {});
    }
    if ('short' in from) {
      to.short = ContentString.copy(from.short || {}, to.short || {});
    }
    if ('text' in from) {
      to.text = ContentText.copy(from.text || {}, to.text || {});
    }
    return to;
  }
}

export abstract class ThemeAsset {
  uuid?: string;
  name?: string;
  theme?: string;
  data?: string;
  static copy(from: ThemeAsset, to?: ThemeAsset): ThemeAsset {
    to = to || {};
    to.uuid = from.uuid;
    to.name = from.name;
    to.theme = from.theme;
    to.data = from.data;
    return to;
  }
}

export abstract class ContentMultiple {
  title?: string;
  options?: string[];
  type?: ContentMultiple_DropdownType;
  static copy(from: ContentMultiple, to?: ContentMultiple): ContentMultiple {
    to = to || {};
    to.title = from.title;
    to.options = from.options;
    to.type = from.type;
    return to;
  }
}

export abstract class ContentText {
  title?: string;
  type?: ContentTextType;
  static copy(from: ContentText, to?: ContentText): ContentText {
    to = to || {};
    to.title = from.title;
    to.type = from.type;
    return to;
  }
}

export abstract class ContentString {
  title?: string;
  type?: ContentTextType;
  static copy(from: ContentString, to?: ContentString): ContentString {
    to = to || {};
    to.title = from.title;
    to.type = from.type;
    return to;
  }
}

export abstract class Author {
  uuid?: string;
  static copy(from: Author, to?: Author): Author {
    to = to || {};
    to.uuid = from.uuid;
    return to;
  }
}

export abstract class Data {
  uuid?: string;
  key?: string;
  value?: string;
  timestamps?: Timestamp;
  // skipped field: type
  
  // oneof types:
  multiple?: ContentMultiple;
  short?: ContentString;
  text?: ContentText;
  static copy(from: Data, to?: Data): Data {
    to = to || {};
    to.uuid = from.uuid;
    to.key = from.key;
    to.value = from.value;
    if ('timestamps' in from) {
      to.timestamps = Timestamp.copy(from.timestamps || {}, to.timestamps || {});
    }
    if ('multiple' in from) {
      to.multiple = ContentMultiple.copy(from.multiple || {}, to.multiple || {});
    }
    if ('short' in from) {
      to.short = ContentString.copy(from.short || {}, to.short || {});
    }
    if ('text' in from) {
      to.text = ContentText.copy(from.text || {}, to.text || {});
    }
    return to;
  }
}

export abstract class Package {
  type?: PackageType;
  name?: string;
  authors?: PackageAuthor[];
  description?: string;
  homepage?: string;
  tags?: string[];
  readmeUrl?: string;
  vcsUrl?: string;
  screenshotUrls?: string[];
  ketchupVersion?: string;
  static copy(from: Package, to?: Package): Package {
    to = to || {};
    to.type = from.type;
    to.name = from.name;
    to.authors = from.authors;
    to.description = from.description;
    to.homepage = from.homepage;
    to.tags = from.tags;
    to.readmeUrl = from.readmeUrl;
    to.vcsUrl = from.vcsUrl;
    to.screenshotUrls = from.screenshotUrls;
    to.ketchupVersion = from.ketchupVersion;
    return to;
  }
}

export abstract class PackageAuthor {
  name?: string;
  email?: string;
  static copy(from: PackageAuthor, to?: PackageAuthor): PackageAuthor {
    to = to || {};
    to.name = from.name;
    to.email = from.email;
    return to;
  }
}

export abstract class Registry {
  registryVersion?: string;
  registryType?: string;
  packages?: Package[];
  static copy(from: Registry, to?: Registry): Registry {
    to = to || {};
    to.registryVersion = from.registryVersion;
    to.registryType = from.registryType;
    to.packages = from.packages;
    return to;
  }
}

export abstract class TLSSettingsResponse {
  tlsEmail?: string;
  tlsDomain?: string;
  agreedOn?: string;
  termsOfService?: string;
  hasCertificate?: boolean;
  static copy(from: TLSSettingsResponse, to?: TLSSettingsResponse): TLSSettingsResponse {
    to = to || {};
    to.tlsEmail = from.tlsEmail;
    to.tlsDomain = from.tlsDomain;
    to.agreedOn = from.agreedOn;
    to.termsOfService = from.termsOfService;
    to.hasCertificate = from.hasCertificate;
    return to;
  }
}

export abstract class EnableTLSRequest {
  tlsEmail?: string;
  tlsDomain?: string;
  agreed?: boolean;
  static copy(from: EnableTLSRequest, to?: EnableTLSRequest): EnableTLSRequest {
    to = to || {};
    to.tlsEmail = from.tlsEmail;
    to.tlsDomain = from.tlsDomain;
    to.agreed = from.agreed;
    return to;
  }
}

export abstract class ListPageRequest {
  list?: ListOptions;
  options?: ListPageRequest_ListPageOptions;
  static copy(from: ListPageRequest, to?: ListPageRequest): ListPageRequest {
    to = to || {};
    if ('list' in from) {
      to.list = ListOptions.copy(from.list || {}, to.list || {});
    }
    if ('options' in from) {
      to.options = ListPageRequest_ListPageOptions.copy(from.options || {}, to.options || {});
    }
    return to;
  }
}

export abstract class ListPageRequest_ListPageOptions {
  filter?: ListPageRequest_ListPageFilter;
  static copy(from: ListPageRequest_ListPageOptions, to?: ListPageRequest_ListPageOptions): ListPageRequest_ListPageOptions {
    to = to || {};
    to.filter = from.filter;
    return to;
  }
}

export abstract class ListPageResponse {
  pages?: Page[];
  static copy(from: ListPageResponse, to?: ListPageResponse): ListPageResponse {
    to = to || {};
    to.pages = from.pages;
    return to;
  }
}

export abstract class ListOptions {
  static copy(_: ListOptions, to?: ListOptions): ListOptions {
    to = to || {};
    return to;
  }
}

export abstract class ListDataResponse {
  data?: Data[];
  static copy(from: ListDataResponse, to?: ListDataResponse): ListDataResponse {
    to = to || {};
    to.data = from.data;
    return to;
  }
}

export abstract class UpdateDataRequest {
  data?: Data[];
  static copy(from: UpdateDataRequest, to?: UpdateDataRequest): UpdateDataRequest {
    to = to || {};
    to.data = from.data;
    return to;
  }
}

export abstract class Error {
  code?: string;
  title?: string;
  detail?: string;
  field?: string;
  static copy(from: Error, to?: Error): Error {
    to = to || {};
    to.code = from.code;
    to.title = from.title;
    to.detail = from.detail;
    to.field = from.field;
    return to;
  }
}

export abstract class ErrorResponse {
  errors?: Error[];
  static copy(from: ErrorResponse, to?: ErrorResponse): ErrorResponse {
    to = to || {};
    to.errors = from.errors;
    return to;
  }
}

export type ContentMultiple_DropdownType = 'dropdown' | 'radio' | 'unknown';
export type ContentTextType = 'html' | 'markdown' | 'text';
export type ListPageRequest_ListPageFilter = 'all' | 'draft' | 'published';
export type PackageType = 'plugin' | 'theme' | 'unknown';
