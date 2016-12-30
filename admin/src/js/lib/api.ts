// DO NOT EDIT! This file is generated automatically by util/gots/main.go

export abstract class Page {
  uuid?: string;
  name?: string;
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
    to.name = from.name;
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
  contentType?: Content_ContentType;
  key?: string;
  value?: string;
  timestamps?: Timestamp;
  static copy(from: Content, to?: Content): Content {
    to = to || {};
    to.uuid = from.uuid;
    to.contentType = from.contentType;
    to.key = from.key;
    to.value = from.value;
    to.timestamps = from.timestamps;
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
  placeholders?: ThemePlaceholder[];
  data?: string;
  static copy(from: ThemeTemplate, to?: ThemeTemplate): ThemeTemplate {
    to = to || {};
    to.uuid = from.uuid;
    to.name = from.name;
    to.theme = from.theme;
    to.engine = from.engine;
    to.placeholders = from.placeholders;
    to.data = from.data;
    return to;
  }
}

export abstract class ThemePlaceholder {
  key?: string;
  type?: ThemePlaceholder_ThemePlaceholderType;
  contentType?: Content_ContentType;
  static copy(from: ThemePlaceholder, to?: ThemePlaceholder): ThemePlaceholder {
    to = to || {};
    to.key = from.key;
    to.type = from.type;
    to.contentType = from.contentType;
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

export type Content_ContentType = 'unknown' | 'html' | 'markdown';
export type ThemePlaceholder_ThemePlaceholderType = 'unknown' | 'string' | 'text';
