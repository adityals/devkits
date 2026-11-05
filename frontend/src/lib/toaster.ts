import { createToaster } from "@skeletonlabs/skeleton-svelte";

// One app-wide toaster; Shell renders it, any module can call it.
export const toaster = createToaster({
  placement: "bottom-end",
});
