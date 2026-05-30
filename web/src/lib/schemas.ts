import { z } from "zod";

export const ErrorResponse = z
    .object({
        code: z.int(),
        request_id: z.string(),
        message: z.string(),
        errors: z
            .array(
                z.object({
                    field: z.string(),
                    message: z.string(),
                }),
            )
            .optional(),
    })
    .transform((d) => ({ ...d, requestId: d.request_id }));

export type ErrorResponse = z.infer<typeof ErrorResponse>;
