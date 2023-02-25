import React, {useEffect, useState} from "react";

const CopyLinkButton = () => {

    const [copied, setCopied] = useState(false);

    const copy = async () => {
        await navigator.clipboard.writeText(window.location.href);
        setCopied(true);
    };

    useEffect(() => {
        if (copied) {
            const timeout = setTimeout(() => {
                setCopied(false);
            }, 3000);
            return () => clearTimeout(timeout);
        }
    }, [copied]);

    return <button className={"bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-800 rounded-lg p-2 cursor-pointer"} onClick={copy}>
        {copied ? "✅ Copied!" : "🔗 Copy link"}
    </button>
}

export default CopyLinkButton;