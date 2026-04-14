import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const res = await data.apiClient.getApiTokens();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    ...data,

    tokens: res.data.tokens,
  };
};
