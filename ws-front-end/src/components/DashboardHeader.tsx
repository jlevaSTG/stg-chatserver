function DashBoardHeader() {
    return <div className="px-4 py-10 sm:px-6 lg:px-8">
        <div className="mx-auto flex items-center justify-between gap-x-8 lg:mx-0">
            <div className="flex items-center gap-x-6">
                <img
                    src="https://tailwindui.com/img/logos/48x48/tuple.svg"
                    alt=""
                    className="h-16 w-16 flex-none rounded-full ring-1 ring-gray-900/10"
                />
                <h1>
                    <div className="mt-1 text-base leading-6 text-gray-700">Stg - Cargo Manager Chat Admin</div>
                </h1>
            </div>

        </div>
    </div>;
}

export default DashBoardHeader;