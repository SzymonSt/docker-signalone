export interface PageDTO<T> {
  number: number;
  numberOfElements: number;
  totalPages: number;
  totalElements: number;
  content: T[];
}