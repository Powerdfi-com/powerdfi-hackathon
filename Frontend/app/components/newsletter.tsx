const NewsLetter = () => {
    return (
        <section className="flex justify-center w-full px-8 sm:px-12 md:px-24 my-16">
            <section className="w-full flex justify-center rounded-2xl bg-gradient-to-br from-white/10 to-black-shade/20 border-b border-primary">
                <article className="py-20 flex flex-col items-center gap-4 max-w-md">
                    <h3 className="text-white font-semibold text-2xl leading-relaxed">Stay up to date</h3>
                    <p className="text-white text-md leading-relaxed text-center">Join our mailing list to stay in the loop with our newest feature releases, tips and tricks for navigating PowerDfi.</p>
                    <form className="ring-1 ring-secondary rounded-lg w-full flex mt-4">
                        <input type="email" placeholder="Your email" className="border-none bg-transparent flex-1 text-sm px-2 py-3" />
                        <button className="bg-secondary rounded-lg px-4 py-1 text-sm">Subscribe</button>
                    </form>
                </article>
            </section>
        </section>
    )
}

export default NewsLetter