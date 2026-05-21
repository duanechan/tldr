import { z } from "zod";

export const ErrorResponseSchema = z.object({
  code: z.number(),
  requestId: z.string(),
  message: z.string(),
  errors: z.array(
    z
      .object({
        field: z.string(),
        message: z.string(),
      })
      .optional(),
  ),
});
