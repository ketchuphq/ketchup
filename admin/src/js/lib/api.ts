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
  multiple?: ContentMultiple;
  short?: ContentString;
  text?: ContentText;
  static copy(from: ThemePlaceholder, to?: ThemePlaceholder): ThemePlaceholder {
    to = to || {};
    to.key = from.key;
    to.multiple = from.multiple;
    to.short = from.short;
    to.text = from.text;
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

export type ContentMultiple_DropdownType = 'unknown' | 'radio' | 'dropdown';
export type ContentTextType = 'text' | 'markdown' | 'html';
