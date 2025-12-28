export interface PaginatedRequest {
  pageIndex: number;
  pageSize: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  pageIndex: number;
  pageSize: number;
  totalCount: number;
  totalPages: number;
}
