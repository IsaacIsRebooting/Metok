import {
  generateSchemaTypes,
  generateReactQueryComponents,
} from "@openapi-codegen/typescript";
import { defineConfig } from "@openapi-codegen/cli";
export default defineConfig({
  svapi: {
    from: {
      relativePath: "backend/shortVideoApiService/openapi.yaml",
      source: "file",
    },
    outputDir: "frontend/metok/api/svapi",
    to: async (context) => {
      const filenamePrefix = "svapi";
      const { schemasFiles } = await generateSchemaTypes(context, {
        filenamePrefix,
      });
      await generateReactQueryComponents(context, {
        filenamePrefix,
        schemasFiles,
      });
    },
  },
});
