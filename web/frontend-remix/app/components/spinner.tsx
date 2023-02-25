import React from "react";

const Spinner = (props: {className: string}) => {
    return (
        <div className={`border-4 rounded-full spinner ${props.className}`} />
    )
}

export default Spinner;