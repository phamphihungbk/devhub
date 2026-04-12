export interface IUserBEResponse {
  id: string
  name: string
  email: string
  role?: string
}

export interface IUser {
  id: IUserBEResponse['id']
  name: IUserBEResponse['name']
  email: IUserBEResponse['email']
  role: IUserBEResponse['role']
}

export class User implements IUser {
  id: IUser['id']
  name: IUser['name']
  email: IUser['email']
  role: IUser['role']

  constructor(data: IUserBEResponse) {
    this.id = data.id
    this.name = data.name
    this.email = data.email
    this.role = data.role
  }
}

export interface FindAllUsersQuery {
  startDate?: string
  endDate?: string
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface Pagination {
  limit?: number
  offset?: number
  total?: number
}

export interface PaginationMetadata {
  pagination?: Pagination
}

export interface CreateUserInput {
  name?: string
  email: string
  password: string
  role: string
}

export interface UpdateUserInput {
  name: string
  role: string
}