import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query';
import { defineBoot } from '#q-app/wrappers';
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
});

export default defineBoot(({ app }) => {
  app.use(VueQueryPlugin, { queryClient });
});

export { queryClient };
