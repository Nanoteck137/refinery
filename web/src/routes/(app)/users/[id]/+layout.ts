import { error } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ parent, params }) => {
  const data = await parent();

  const user = await data.apiClient.getUser(params.id);
  if (!user.success) {
    return error(user.error.code, { message: user.error.message });
  }

  return {
    ...data,

    userData: user.data.user,
  };
};
