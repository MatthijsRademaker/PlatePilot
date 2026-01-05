import { defineConfig } from 'orval';

export default defineConfig({
  platepilot: {
    output: {
      mode: 'split',
      target: 'src/api/generated/platepilot.ts',
      schemas: 'src/api/generated/models',
      client: 'vue-query',
      mock: false,
      override: {
        mutator: {
          path: './src/api/mutator/custom-instance.ts',
          name: 'customInstance',
        },
        query: {
          useQuery: true,
          useMutation: true,
        },
      },
    },
    input: {
      target: '../backend/api/openapi/swagger.yaml',
    },
  },
});
