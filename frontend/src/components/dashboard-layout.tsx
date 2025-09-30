import type { PropsWithChildren } from 'react';
import type { KeyedMutator } from 'swr';
import type z from 'zod';
import { Sidebar } from '@components/sidebar/sidebar';
import { zodResolver } from '@hookform/resolvers/zod';
import { DialogContent, DialogDescription, DialogOverlay, DialogPortal, DialogTitle, DialogTrigger, Root } from '@radix-ui/react-dialog';
import { Button } from '@ui/button';
import { FormField, Input, Textarea } from '@ui/form';
import clsx from 'clsx';
import { Loader2Icon, RefreshCcwIcon } from 'lucide-react';
import { useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { testTriggerSchema } from '@/lib/schemas';

export function DashoardLayout({ title, mutate, children }: { title: string; mutate: KeyedMutator<any> } & PropsWithChildren) {
  const [isMutating, setIsMutating] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm({
    resolver: zodResolver(testTriggerSchema),
    defaultValues: {
      payload: '{}',
      function_name: '',
    },
  });

  async function handleMutate() {
    setIsMutating(true);
    await mutate();
    setIsMutating(false);
  }

  async function onSubmit(values: z.infer<typeof testTriggerSchema>) {
    setIsLoading(true);
    const result = await fetch(`${import.meta.env.VITE_API_URL}/trigger`, {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        function_name: values.function_name,
        payload: values.payload ? JSON.stringify(values.payload) : undefined,
      }),
    });
    setIsLoading(false);
    if (!result?.ok) {
      const error = await result.json();
      toast.error('Failed to trigger', {
        description: error?.message ?? 'Unknown error.',
      });
      return;
    }
    toast.success('Triggered successfully');
    form.reset();
    mutate();
  }

  return (
    <div className="flex w-full">
      <Sidebar />
      <div className="flex flex-col w-full">
        <div className="p-6 border-b border-neutral-800 flex items-center justify-between h-18 gap-2">
          <span className="text-lg font-semibold">{title}</span>
          <div className="flex items-center gap-2">
            <Root>
              <DialogTrigger asChild>
                <Button size="sm">Send test trigger</Button>
              </DialogTrigger>
              <DialogPortal>
                <DialogOverlay className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
                <DialogContent className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl bg-neutral-950 p-6 shadow-sm border border-neutral-800 text-white">
                  <DialogTitle className="text-xl font-semibold">Send test trigger</DialogTitle>
                  <DialogDescription className="text-sm mt-0.5 text-neutral-400">You can send a test trigger to the server.</DialogDescription>
                  <div className="mt-6">
                    <FormProvider {...form}>
                      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                          label="Function name"
                          name="function_name"
                          control={form.control}
                          render={field => (
                            <Input {...field} type="text" placeholder="function-name" />
                          )}
                        />
                        <FormField
                          label="Event payload"
                          name="payload"
                          control={form.control}
                          render={field => (
                            <Textarea {...field} type="text" placeholder={`{"message": "hello, world"}`} />
                          )}
                        />
                        <div className="mt-6">
                          <Button type="submit" disabled={isLoading}>
                            {isLoading && <Loader2Icon className="animate-spin" size="16" />}
                            Trigger
                          </Button>
                        </div>
                      </form>
                    </FormProvider>
                  </div>
                </DialogContent>
              </DialogPortal>
            </Root>
            <Button size="icon" onClick={handleMutate} disabled={isMutating}>
              <RefreshCcwIcon size={16} className={clsx(isMutating && 'animate-spin')} />
            </Button>
          </div>
        </div>
        <div className="flex-1 p-6">
          {children}
        </div>
      </div>
    </div>
  );
}
