/**
 * User is an account holder.
 */
export interface User {
  id: number /* int64 */;
  email?: string;
  tags?: string[];
  meta: { [key: string]: string};
  created_at: string;
  address: Address;
}
export interface Address {
  city: string;
}
