"use client"
import Image from "next/image";

export default function AboutUs() {
  return (
    <main className="bg-gradient-radial from-black to-primary/5">
      <section className="flex flex-col justify-center items-center py-24  w-full">
        <div className="absolute top-40 bottom-0 right-20 -z-10">
          <Image src="/logo_icon.png" height={642} width={638} alt="logo" className="flex justify-center items-center" />
        </div>
        <div className="flex justify-center items-center text-[64px] font-bold text-white text-center max-w-[1069px]">
          <div>
            <span className="text-white ">About  </span><span className="text-secondary">PowerDfi </span>
          </div>
          <div className="w-[146px] h-[146px] relative">
            <Image src="/icon.png" fill={true} alt="logo" className="flex justify-center items-center" />
          </div>

        </div>
        <div className="max-w-[578px] pb-16 border-b border-primary/40">
          <p className="text-white/60 text-[24px] leading-relaxed text-center mt-6">
            We aim to revolutionise the way real-world assets (RWAs) are managed and traded by leveraging blockchain technology.
          </p>
        </div>
      </section>

      <section className="px-24 pt-4 pb-16 w-full">
        <div className="flex justify-center items-center">
          <div className="flex-1 max-w-[635px]">
            <h3 className="font-semibold leading-relaxed text-secondary text-[36px] max-w-[522px]">Revolutionizing Real-World Assets with PowerDfi</h3>
            <p className="text-[24px] leading-relaxed text-white mt-8 max-w-[582px]">PowerDfi was established to create a fairer financial landscape for individuals,
              families, communities, and nations. Our founder recognised the significant inequities in the global financial system and envisioned using blockchain
              technology to streamline financial services, enabling more people to achieve economic freedom.</p>
          </div>

          <div className="mt-20">
            <div className="relative w-96">
              <div className="absolute -left-5 -top-5 bg-purple-700 h-full w-full rounded-xl"></div>
              <div className="relative">
                <Image src={"/asset.png"} height={396} width={490} alt="info" className="rounded-lg" />
              </div>
            </div>

          </div>

        </div>
      </section>
      <section className="py-16 px-8 sm:px-12 md:px-24 w-full bg-gradient-to-r from-black to-gray-700">
        {/*<div className="bg-[url(/logo_icon.png)] w-[500px] h-[500px] bg-center bg-cover">*/}
        <div className="flex justify-center items-center gap-12 border-b border-primary/40">

          <div className="relative flex-1 h-[600px]">
            <Image src={"/asset-2.png"} alt="" fill={true} className="rounded-xl relative" />
          </div>
          <div className="flex flex-1 flex-col text-secondary text-[128px] font-bold gap-0 ">
            <span>Who</span><span className="ml-28 -mt-20">We</span><span className="ml-56 -mt-20">Are</span>
          </div>
        </div>

        <div className="absolute flex justify-center items-center -mt-96 w-full h-full">
          <Image src="/logo_icon.png" height={1024} width={1024} alt="logo" className="" />
        </div>

        <div className="py-8 ml-16 flex justify-center  gap-24">
          <div>
            <h3 className="font-bold leading-relaxed text-secondary text-[36px] ">Our Vision</h3>
            <p className="text-[24px] leading-relaxed text-white mt-4 max-w-[582px]">To unlock trillions of frozen assets globally, transforming them
              into accessible and liquid digital tokens.</p>
          </div>
          <div>
            <h3 className="font-bold leading-relaxed text-secondary text-[36px] ">Our Mision</h3>
            <p className="text-[24px] leading-relaxed text-white mt-4 max-w-[582px]">To level the financial playing field for individuals, families,
              communities, and nations by leveraging the latest fintech technologies, such as blockchain, for a more inclusive and equitable financial system.</p>
          </div>

        </div>

      </section>

      <section className="flex flex-col justify-center w-full px-8 sm:px-12 md:px-24 my-16">
        <div className="py-8 flex flex-col justify-center items-center  gap-12 text-[32px] text-white text-center px-8 mb-20">
          <p className="max-w-[1024px]">
            At PowerDfi, we transform Real World Assets into digital tokens, offering enhanced transparency and liquidity for
            asset owners. This innovation is essential for asset-rich individuals,
            communities, firms, and nations, helping unlock trillions of frozen assets globally.
          </p>
          <p className="max-w-[1024px]">
            Join us in our mission. PowerDfi is not just a service; we are a community of wealth builders for everyone.
          </p>

        </div>
        <section className="w-full flex justify-center rounded-2xl bg-gradient-to-br from-white/10 to-black-shade/20 border-b border-primary">
          <article className="py-20 flex flex-col items-center gap-4 max-w-md">
            <h3 className="text-white font-semibold text-2xl leading-relaxed">Join the Mission</h3>
            <p className="text-white text-md leading-relaxed text-center">Join our mailing list to stay in the loop with our newest feature releases, tips and tricks for navigating PowerDfi.</p>
            <form className="ring-1 ring-secondary rounded-lg w-full flex mt-4">
              <input placeholder="Your email" className="border-none bg-transparent flex-1 text-sm px-2 py-3" type="email"></input>
              <button className="bg-secondary rounded-lg px-4 py-1 text-sm">Subscribe</button>
            </form>
          </article>
        </section>
      </section>
    </main>
  );
}
