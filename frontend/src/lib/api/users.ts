import { apiV1client } from "./axios";

export type User = {
  email: string;
};

export type GetUsersRes = {
  ok: boolean;
  error?: string;
  users: User[];
};
export function getUsers() {
  return apiV1client.get<GetUsersRes>("/api/v1/users");
}
export type CreateUserRes = {
  ok: boolean;
  error?: string;
};
export function createUser(creds: FormData) {
  return apiV1client.post<CreateUserRes>("/api/v1/users", creds);
}

export type DeleteUserRes = {
  ok: boolean;
  error?: string;
};
export function deleteUser(email: string) {
  return apiV1client.delete<DeleteUserRes>("/api/v1/users", {
    params: {
      email,
    },
  });
}

export type UpdateUserPasswordRes = {
  ok: boolean;
  error?: string;
};
export function updateUserPassword(creds: FormData) {
  return apiV1client.patch<UpdateUserPasswordRes>("/api/v1/users", creds);
}
