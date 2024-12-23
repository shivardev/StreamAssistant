import { rootAction } from "./commentApi.modal copy";

export interface CommentObj {
  responseContext: ResponseContext;
  continuationContents: ContinuationContents;
  liveChatStreamingResponseExtension: LiveChatStreamingResponseExtension;
}

interface LiveChatStreamingResponseExtension {
  lastPublishAtUsec: string;
}

interface ContinuationContents {
  liveChatContinuation: LiveChatContinuation;
}

interface LiveChatContinuation {
  actions?: rootAction[];
  continuations: Continuation[];
}

interface Continuation {
  invalidationContinuationData: InvalidationContinuationData;
}

interface InvalidationContinuationData {
  invalidationId: InvalidationId;
  timeoutMs: number;
  continuation: string;
}

interface InvalidationId {
  objectSource: number;
  objectId: string;
  topic: string;
  subscribeToGcmTopics: boolean;
  protoCreationTimestampMs: string;
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
  trackingParam: string;
}

interface ServiceTrackingParam {
  service: string;
  params: Param[];
}

interface Param {
  key: string;
  value: string;
}
