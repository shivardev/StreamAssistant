export interface rootAction {
  addChatItemAction: AddChatItemAction;
}

interface AddChatItemAction {
  item: Item;
  clientId: string;
}

interface Item {
  liveChatTextMessageRenderer: LiveChatTextMessageRenderer;
}

interface LiveChatTextMessageRenderer {
  message: Message;
  authorName: AuthorName;
  authorPhoto: AuthorPhoto;
  contextMenuEndpoint: ContextMenuEndpoint;
  id: string;
  timestampUsec: string;
  authorExternalChannelId: string;
  contextMenuAccessibility: ContextMenuAccessibility;
}

interface ContextMenuAccessibility {
  accessibilityData: AccessibilityData;
}

interface AccessibilityData {
  label: string;
}

interface ContextMenuEndpoint {
  commandMetadata: CommandMetadata;
  liveChatItemContextMenuEndpoint: LiveChatItemContextMenuEndpoint;
}

interface LiveChatItemContextMenuEndpoint {
  params: string;
}

interface CommandMetadata {
  webCommandMetadata: WebCommandMetadata;
}

interface WebCommandMetadata {
  ignoreNavigation: boolean;
}

interface AuthorPhoto {
  thumbnails: Thumbnail[];
}

interface Thumbnail {
  url: string;
  width: number;
  height: number;
}

interface AuthorName {
  simpleText: string;
}

interface Message {
  runs: Run[];
}

interface Run {
  text?: string;
  emoji?: Emoji;
}
interface Emoji {
  emojiId: string;
  shortcuts: string[];
  searchTerms: string[];
  supportsSkinTone: boolean;
  image: Image;
  variantIds: string[];
}
interface Image {
  thumbnails: { url: string }[];
  accessibility: Accessibility;
}
interface Accessibility {
  accessibilityData: AccessibilityData;
}