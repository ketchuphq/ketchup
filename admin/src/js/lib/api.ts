// DO NOT EDIT! This file is generated automatically by util/gots/main.go

export interface Page {
  uuid?: string;
  name?: string;
  theme?: string;
  template?: string;
  timestamps?: Timestamp;
  publishedAt?: string;
  contents?: Content[];
  metadata?: { [key: string]: string; };
  tags?: string[];
}

export interface Content {
  uuid?: string;
  contentType?: Content_ContentType;
  key?: string;
  value?: string;
  timestamps?: Timestamp;
}

export interface Route {
  uuid?: string;
  path?: string;
  // skipped field: target
  
  // oneof types:
  pageUuid?: string;
  file?: string;
}

export interface Timestamp {
  createdAt?: string;
  updatedAt?: string;
}

export interface Theme {
  uuid?: string;
  name?: string;
  templates?: { [key: string]: ThemeTemplate; };
  assets?: { [key: string]: ThemeAsset; };
}

export interface ThemeTemplate {
  uuid?: string;
  name?: string;
  theme?: string;
  engine?: string;
  placeholders?: ThemePlaceholder[];
  data?: string;
}

export interface ThemePlaceholder {
  name?: string;
  type?: ThemePlaceholder_ThemePlaceholderType;
  contentType?: Content_ContentType;
}

export interface ThemeAsset {
  uuid?: string;
  name?: string;
  theme?: string;
  data?: string;
}

export type Content_ContentType = 'unknown' | 'html' | 'markdown';
export type ThemePlaceholder_ThemePlaceholderType = 'unknown' | 'string' | 'text';
