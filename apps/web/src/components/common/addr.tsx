interface AddrOptions {
    leftLength?: number;
    rightLength?: number;
}

export const addr = (address: string, { leftLength = 5, rightLength = 5 }: AddrOptions = {}) => {
    const leftPart = address.slice(0, leftLength);
    const rightPart = address.slice(-rightLength);

    return `${leftPart}...${rightPart}`;
};
