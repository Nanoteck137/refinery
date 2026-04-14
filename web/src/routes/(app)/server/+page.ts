import { isRoleAdmin } from "$lib/utils";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  if (!isRoleAdmin(data.user?.role ?? "")) {
    redirect(301, "/");
  }

  return {
    ...data,
  };
};
