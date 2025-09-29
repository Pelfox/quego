import z from 'zod';

export const testTriggerSchema = z.object({
  function_name: z.string('You must use a valid string.')
    .min(1, 'Function name is required'),
  payload: z.string('You must use a valid string.')
    .optional(),
});
