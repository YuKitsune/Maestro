import {Link} from "remix";
import ArrowIcon, {ArrowDirection} from "~/components/icons/arrow";

const HomeButton =() => {
    return (
        <div className={"flex flex-initial hover:underline"}>
            <Link to={"/"} className={"flex flex-row flex-initial content-center"}>
                <ArrowIcon direction={ArrowDirection.Left} className={"h-6 w-6 pr-2"} />
                <div className={""}>Home</div>
            </Link>
        </div>
    );
}

export default HomeButton;