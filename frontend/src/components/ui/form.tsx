import type { InputHTMLAttributes, PropsWithChildren, ReactNode } from 'react';
import type { Control, ControllerRenderProps, FieldPath, FieldValues } from 'react-hook-form';
import clsx from 'clsx';
import { Controller, useFormContext } from 'react-hook-form';

export function Label({ children }: PropsWithChildren) {
  return (
    <label className="block mb-1 text-sm text-neutral-200">
      {children}
    </label>
  );
}

export function Input({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={clsx(
        'text-sm w-full bg-neutral-900 text-white rounded-lg border border-neutral-800 py-2 px-3 placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-orange-500/50 focus:border-orange-500 transition',
        className,
      )}
    />
  );
}

export function Textarea({ className, ...props }: InputHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      {...props}
      className={clsx(
        'text-sm w-full bg-neutral-900 text-white rounded-lg border border-neutral-800 py-2 px-3 placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-orange-500/50 focus:border-orange-500 transition',
        className,
      )}
    />
  );
}

interface FormFieldProps<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> {
  label: string;
  name: TName;
  control: Control<TFieldValues>;
  render: (field: ControllerRenderProps<TFieldValues, TName>) => ReactNode;
}

export function FormField<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ label, name, control, render }: FormFieldProps<TFieldValues, TName>) {
  const context = useFormContext();
  if (!context)
    throw new Error('FormField must be used within a FormProvider');
  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => (
        <div>
          <Label>{label}</Label>
          {render(field)}
          {context.formState.errors[name] && (
            <p className="text-red-500 text-sm mt-1">{context.formState.errors[name].message?.toString()}</p>
          )}
        </div>
      )}
    />
  );
}
