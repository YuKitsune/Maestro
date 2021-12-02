
export enum ArrowDirection {
    Up,
    Down,
    Left,
    Right
}

type ArrowIconProps = {
    className?: string;
    direction: ArrowDirection;
}

const ArrowIcon = (props: ArrowIconProps) => {
    return <svg xmlns="http://www.w3.org/2000/svg" className={`${directionClassName(props.direction)} ${props.className}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
    </svg>;
}

const directionClassName = (dir: ArrowDirection): string => {
    switch (dir) {
        case ArrowDirection.Up:
            return "rotate-90";
        case ArrowDirection.Down:
            return "rotate-270";
        case ArrowDirection.Left:
            return "rotate-0";
        case ArrowDirection.Right:
            return "rotate-180";
    }
}

export default ArrowIcon;
