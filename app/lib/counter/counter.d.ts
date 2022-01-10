export interface Counter {
    value: number;
}
export declare function getCounter(): Promise<Counter>;
export declare function incrementCounter(): Promise<Counter>;
