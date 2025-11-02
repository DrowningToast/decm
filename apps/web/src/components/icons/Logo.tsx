import type { SVGProps } from "react";
const SvgLogo = (props: SVGProps<SVGSVGElement>) => (
    <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        overflow="visible"
        preserveAspectRatio="none"
        style={{
            display: "block",
        }}
        viewBox="0 0 31 13"
        width="1em"
        height="1em"
        {...props}
    >
        <path
            stroke="var(--stroke-0, #FCFCFC)"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={3}
            d="m14.779 2 5.808 4.64L5.485 2l22.072 9.281L28.72 2 2 6.64"
        />
    </svg>
);
export default SvgLogo;
