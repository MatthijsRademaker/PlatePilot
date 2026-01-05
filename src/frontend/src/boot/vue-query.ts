import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query';
import type { boot } from 'quasar/wrappers';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
});

export default ((({ app }) => {
  app.use(VueQueryPlugin, { queryClient });
}) as unknown) as ReturnType<typeof boot>;

export { queryClient };
