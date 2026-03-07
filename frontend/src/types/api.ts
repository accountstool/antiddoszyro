export type Envelope<T> = {
  success: boolean;
  message?: string;
  data: T;
  pagination?: Pagination;
};

export type Pagination = {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
};

export type User = {
  id: string;
  username: string;
  email: string;
  displayName: string;
  role: string;
  language: string;
  lastLoginAt?: string | null;
};

export type TimePoint = {
  label: string;
  allowed: number;
  blocked: number;
  challenge: number;
};

export type RankedMetric = {
  name: string;
  value: number;
};

export type DashboardSummary = {
  healthy: boolean;
  totalDomains: number;
  blockedToday: number;
  currentRps: number;
  currentBlockedPerSecond: number;
  topAttackedDomain: string;
  topAttackingIp: string;
  totalRequests24h: number;
  allowed24h: number;
  blocked24h: number;
  challenged24h: number;
  challengePassRate: number;
};

export type Domain = {
  id: string;
  name: string;
  originHost: string;
  originPort: number;
  originProtocol: "http" | "https";
  originServerName: string;
  enabled: boolean;
  protectionEnabled: boolean;
  protectionMode: string;
  challengeMode: string;
  cloudflareMode: boolean;
  sslAutoIssue: boolean;
  sslEnabled: boolean;
  forceHttps: boolean;
  rateLimitRps: number;
  rateLimitBurst: number;
  badBotMode: boolean;
  headerValidation: boolean;
  jsChallengeEnabled: boolean;
  allowedMethods: string[];
  notes: string;
  createdAt: string;
  updatedAt: string;
};

export type DomainRule = {
  id?: string;
  name: string;
  type: string;
  pattern: string;
  action: string;
  enabled: boolean;
  priority: number;
};

export type DomainDetail = {
  domain: Domain;
  rules: DomainRule[];
  nginxStatus: string;
};

export type StatsOverview = {
  incomingRequests: number;
  allowedRequests: number;
  blockRequests: number;
  challengedRequests: number;
  challengePassRate: number;
  uniqueIps: number;
  peakRps: number;
  peakTime: string;
  topIps: RankedMetric[];
  topUserAgents: RankedMetric[];
  topDomains: RankedMetric[];
  topReasons: RankedMetric[];
  topCountries: RankedMetric[];
  requestSeries: TimePoint[];
};

export type RequestLog = {
  id: string;
  domainName: string;
  clientIp: string;
  countryCode: string;
  method: string;
  path: string;
  queryString: string;
  userAgent: string;
  decision: string;
  statusCode: number;
  blockReason: string;
  responseTimeMs: number;
  score: number;
  challengeType: string;
  createdAt: string;
};

export type IPEntry = {
  id: string;
  domainId?: string | null;
  listType: "blacklist" | "whitelist";
  ip: string;
  cidr: string;
  reason: string;
  expiresAt?: string | null;
  createdAt: string;
};

export type TemporaryBan = {
  id: string;
  domainId?: string | null;
  ip: string;
  reason: string;
  source: string;
  expiresAt: string;
  createdAt: string;
};

export type SystemSetting = {
  key: string;
  value: string;
  type: string;
};

export type AuditLog = {
  id: string;
  username: string;
  action: string;
  entityType: string;
  entityId: string;
  ipAddress: string;
  userAgent: string;
  details: string;
  createdAt: string;
};
