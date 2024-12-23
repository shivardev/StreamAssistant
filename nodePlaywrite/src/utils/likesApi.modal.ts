export interface LikeResponseObj {
  responseContext: ResponseContext;
  continuation: Continuation;
  actions: Action[];
  frameworkUpdates: FrameworkUpdates;
}

interface FrameworkUpdates {
  entityBatchUpdate: EntityBatchUpdate;
}

interface EntityBatchUpdate {
  mutations: Mutation[];
  timestamp: Timestamp;
}

interface Timestamp {
  seconds: string;
  nanos: number;
}

interface Mutation {
  entityKey: string;
  type: string;
  payload: Payload;
}

interface Payload {
  likeCountEntity: LikeCountEntity;
}

interface LikeCountEntity {
  key: string;
  likeCountIfLiked: LikeCountIfLiked;
  likeCountIfDisliked: LikeCountIfLiked;
  likeCountIfIndifferent: LikeCountIfLiked;
  expandedLikeCountIfLiked: LikeCountIfLiked;
  expandedLikeCountIfDisliked: LikeCountIfLiked;
  expandedLikeCountIfIndifferent: LikeCountIfLiked;
  likeCountLabel: LikeCountIfLiked;
  likeButtonA11yText: LikeCountIfLiked;
  likeCountIfLikedNumber: string;
  likeCountIfDislikedNumber: string;
  likeCountIfIndifferentNumber: string;
  shouldExpandLikeCount: boolean;
  sentimentFactoidA11yTextIfLiked: LikeCountIfLiked;
  sentimentFactoidA11yTextIfDisliked: LikeCountIfLiked;
}

interface LikeCountIfLiked {
  content: string;
}

interface Action {
  updateViewershipAction: UpdateViewershipAction;
}

interface UpdateViewershipAction {
  viewCount: ViewCount2;
}

interface ViewCount2 {
  videoViewCountRenderer: VideoViewCountRenderer;
}

interface VideoViewCountRenderer {
  viewCount: ViewCount;
  isLive: boolean;
  extraShortViewCount: ExtraShortViewCount;
  unlabeledViewCountValue: ViewCount;
  originalViewCount: string;
}

interface ExtraShortViewCount {
  accessibility: Accessibility;
  simpleText: string;
}

interface Accessibility {
  accessibilityData: AccessibilityData;
}

interface AccessibilityData {
  label: string;
}

interface ViewCount {
  simpleText: string;
}

interface Continuation {
  timedContinuationData: TimedContinuationData;
}

interface TimedContinuationData {
  timeoutMs: number;
  continuation: string;
}

interface ResponseContext {
  serviceTrackingParams: ServiceTrackingParam[];
  mainAppWebResponseContext: MainAppWebResponseContext;
  webResponseContextExtensionData: WebResponseContextExtensionData;
}

interface WebResponseContextExtensionData {
  hasDecorated: boolean;
}

interface MainAppWebResponseContext {
  loggedOut: boolean;
}

interface ServiceTrackingParam {
  service: string;
  params: Param[];
}

interface Param {
  key: string;
  value: string;
}
