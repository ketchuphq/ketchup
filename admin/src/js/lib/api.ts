// DO NOT EDIT! This file is generated automatically by util/gots/main.go

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
  static copy(from: Page, to?: Page): Page {
    to = to || {};
    to.uuid = from.uuid;
    to.title = from.title;
    to.theme = from.theme;
    to.template = from.template;
    to.timestamps = from.timestamps;
    to.publishedAt = from.publishedAt;
    to.contents = from.contents;
    to.metadata = from.metadata;
    to.tags = from.tags;
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
  short?: ContentString;
  text?: ContentText;
  multiple?: ContentMultiple;
  static copy(from: Content, to?: Content): Content {
    to = to || {};
    to.uuid = from.uuid;
    to.key = from.key;
    to.value = from.value;
    to.timestamps = from.timestamps;
    to.short = from.short;
    to.text = from.text;
    to.multiple = from.multiple;
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
  templates?: { [key: string]: ThemeTemplate; };
  assets?: { [key: string]: ThemeAsset; };
  static copy(from: Theme, to?: Theme): Theme {
    to = to || {};
    to.uuid = from.uuid;
    to.name = from.name;
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
  placeholders?: ThemePlaceholder[];
  data?: string;
  static copy(from: ThemeTemplate, to?: ThemeTemplate): ThemeTemplate {
    to = to || {};
    to.uuid = from.uuid;
    to.name = from.name;
    to.theme = from.theme;
    to.engine = from.engine;
    to.hideContent = from.hideContent;
    to.placeholders = from.placeholders;
    to.data = from.data;
    return to;
  }
}

export abstract class ThemePlaceholder {
  key?: string;
  // skipped field: type
  
  // oneof types:
  short?: ContentString;
  text?: ContentText;
  multiple?: ContentMultiple;
  static copy(from: ThemePlaceholder, to?: ThemePlaceholder): ThemePlaceholder {
    to = to || {};
    to.key = from.key;
    to.short = from.short;
    to.text = from.text;
    to.multiple = from.multiple;
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

export abstract class Package {
  name?: string;
  author?: string[];
  type?: Package_Type;
  description?: string;
  readme?: string;
  vcsUrl?: string;
  releases?: PackageRelease;
  labels?: string[];
  static copy(from: Package, to?: Package): Package {
    to = to || {};
    to.name = from.name;
    to.author = from.author;
    to.type = from.type;
    to.description = from.description;
    to.readme = from.readme;
    to.vcsUrl = from.vcsUrl;
    to.releases = from.releases;
    to.labels = from.labels;
    return to;
  }
}

export abstract class PackageRelease {
  ketchupMin?: string;
  tags?: boolean;
  static copy(from: PackageRelease, to?: PackageRelease): PackageRelease {
    to = to || {};
    to.ketchupMin = from.ketchupMin;
    to.tags = from.tags;
    return to;
  }
}

export abstract class Registry {
  registryVersion?: string;
  packages?: Package[];
  static copy(from: Registry, to?: Registry): Registry {
    to = to || {};
    to.registryVersion = from.registryVersion;
    to.packages = from.packages;
    return to;
  }
}

export abstract class TLSSettingsReponse {
  tlsEmail?: string;
  tlsDomain?: string;
  agreedOn?: string;
  termsOfService?: string;
  hasCertificate?: boolean;
  static copy(from: TLSSettingsReponse, to?: TLSSettingsReponse): TLSSettingsReponse {
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

export type ContentMultiple_DropdownType = 'unknown' | 'radio' | 'dropdown';
export type ContentTextType = 'text' | 'markdown' | 'html';
export type Package_Type = 'unknown' | 'theme' | 'plugin';
