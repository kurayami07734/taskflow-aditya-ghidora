declare module 'lodash.debounce' {
  function debounce<T extends (...args: string[]) => unknown>(
    func: T,
    wait?: number
  ): T;
  export = debounce;
}