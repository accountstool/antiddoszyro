export function Spinner() {
  return (
    <div className="relative h-11 w-11">
      <div className="absolute inset-0 animate-spin rounded-full border-[3px] border-slate-200/80 border-t-black dark:border-slate-800 dark:border-t-white" />
      <div className="absolute inset-[7px] rounded-full bg-white/70 dark:bg-neutral-950/70" />
    </div>
  );
}
