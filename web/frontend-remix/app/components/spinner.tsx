import React from "react";

type SpinnerProps = {
    className?: string;
}

const Spinner = (props: SpinnerProps) => {
    return (
        <div className="w-12 h-12 border-4 border-blue-500 rounded-full spinner" />
    )
}

export default Spinner;