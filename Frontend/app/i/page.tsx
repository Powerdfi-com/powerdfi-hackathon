"use client"
import Image from "next/image";
import CustomSlider from "@/app/components/slider";
import { motion } from 'framer-motion';
import { Link } from "@/app/components/link";
import { useAccount } from "wagmi";

import { FaRegEdit } from "react-icons/fa";
import { useContext } from "react";
import { ThemeContext } from "../context/context";
import { BsLightningCharge } from "react-icons/bs";

export default function Home() {
  const { isConnected } = useAccount()
  const infos = [
    {
      image: "/info1.png",
      title: "Store Asset Value, Leverage/Monetise, Trade, and Barter on your terms."
    },
    {
      image: "/info2.png",
      title: "The PowerDfi platform makes it easy for individuals, firms and government in a secure environment."
    },
  ];
  const { isDarkMode } = useContext(ThemeContext);
  return (
    <main className="bg-white dark:bg-black dark:bg-gradient-radial dark:from-black dark:to-primary/5">
      <section className="flex flex-col justify-center items-center py-24 w-full">
        <div className="absolute top-0 bottom-0 right-0 left-0 -z-10">
          <div className="w-full flex h-full items-center justify-center ">
            <div className="h-[565px] w-[565px] gradient rounded-full ring-1 ring-primary/10"></div>
          </div>
        </div>
        <div className="text-[64px] font-semibold text-black dark:text-white text-center max-w-[1069px] leading-[70px]">
          {/* <BsLightningCharge className="text-secondary" /> */}
          <div className="flex gap-4 items-center justify-center">
            <span>We are building the </span><BsLightningCharge className="!text-7xl text-secondary font-bold" /><span className="text-secondary">Future </span><span> of</span>
          </div><br></br><span className="text-secondary"> Real World Asset (RWA) </span><span>Digital banking</span>
        </div>
        <p className="text-black dark:text-white text-[30px] leading-relaxed text-center mt-6">
          Tokenise/Digitise your Real World Assets on our innovative platform
        </p>
        <div className="flex mt-24">
          <button className="text-secondary ring-1 ring-primary px-16 py-2 rounded-sm flex items-center justify-center gap-2 text-[24px]"><FaRegEdit className="!text-xl" /><span>More Information</span></button>
          <Link href={isConnected ? "/u/h" : "/auth"} className="bg-secondary  px-16 py-4 rounded-sm text-[24px]">Join PowerDfi Studio</Link>
        </div>
      </section>
      <section className="flex justify-center my-12 ">
        <section className="max-w-[982px]">
          <CustomSlider>
            {
              infos.map((info) => <div key={info.title} className="my-6 rounded-2xl ring-1 ring-primary/40 !flex dark:bg-transparent dark:gradient bg-text-blue h-[384px] gap-16 max-w-[980px]">
                <div className="relative aspect-square h-[344px] w-[344px]">
                  <Image src={info.image} alt={"info"} fill={true} />
                </div>
                <p className="flex-[2] text-[40px] px-6 leading-relaxed text-black dark:text-white mt-8">
                  {info.title}
                </p>
              </div>)
            }
          </CustomSlider>
        </section>
      </section>
      <ul className="h-14 w-full flex justify-evenly items-center bg-text-blue dark:bg-black-shade">
        {
          ["$1B  Total Value of Tokenized US Treasuries", "$12B Total Loan Value of Tokenized  Private Credit", "$159B Total Value of ‘Stablecoins’"].map((str) => <li key={str}><div className="h-[40px] flex items-center justify-center italic px-2 rounded-full ring-1 ring-text-grey text-black dark:text-white text-[20px] leading-relaxed " >{str}</div></li>)
        }
      </ul>
      <section className="flex justify-center items-center my-12 flex-col w-full px-24">
        <div className="flex w-full justify-start">
          <div className="py-2 px-4 text-[40px] rounded-md ring-1 ring-primary/40 text-secondary font-semibold gradient">The Challenge</div>
        </div>
        <div className="grid grid-cols-3 justify-center gap-12 mt-24 items-center">
          {[{ img: "/challenge1.png", title: "Access To Capital & Financing Is Limited" }, { img: "/challenge2.png", title: "Liquidity For Certain Assets Is Non-Existent" }, { img: "/challenge3.png", title: "Current Banking System Is Inflexible" }, { img: "/challenge4.png", title: "Asset Maintenance Costs Are Exorbitant" }].map(({ img, title }) => <div key={title} className="relative">
            <div className="absolute h-[247px] w-[350px] ring-1 ring-primary/20 rounded-lg -top-4 -left-4"></div>
            <div key={title} className="w-[350px] h-[247px] bg-white dark:gradient ring-1 ring-primary/20 rounded-lg flex flex-col justify-center items-center gap-2 relative">
              <div className="relative h-[100px] w-[100px]">
                <Image src={img} alt="challenge" fill={true} />
              </div>
              <div className="text-[32px] px-4 text-black dark:text-white text-center">{title}</div>
            </div>
          </div>)}
        </div>
      </section>
      <section id="features" className="flex justify-center items-center my-24 flex-col w-full px-24">
        <div className="flex w-full justify-start">
          <div className="py-2 px-4 text-[40px] rounded-md ring-1 ring-primary/40 text-secondary font-semibold ]">Our Solution</div>
        </div>
        <div className="flex gap-12 items-center my-6 w-full max-w-[940px] mt-24">
          <div className="relative h-[256px] w-[340px] max-w-xs">
            <div className="h-[256px] w-[340px] bg-gradient-to-br from-primary/10 to-primary/5 rounded-lg ring-1 ring-primary/20 absolute bottom-6 right-6">
            </div>
            <Image src={"/lux.png"} fill={true} alt="info" className="rounded-lg" />
          </div>
          <p className="max-w-[407px] text-black dark:text-white text-[40px] leading-relaxed">Tokenise/Digitise your Real World Asset</p>
        </div>
        <div className="flex gap-12 items-center justify-end my-6 w-full max-w-[940px]">
          <p className="max-w-[407px] text-[40px] text-black dark:text-white leading-relaxed text-right">Trade/Barter your
            Real Word Asset</p>
          <div className="relative h-[256px] w-[340px]  flex-1 max-w-xs">
            <div className="bg-gradient-to-br from-primary/10 to-primary/5 rounded-lg ring-1 ring-primary/20 absolute bottom-6 left-6 h-[256px] w-[340px] ">
            </div>
            <Image src={"/painting.png"} fill={true} alt="info" className="rounded-lg" />
          </div>
        </div>
        <div className="flex gap-12 items-center my-6 w-full max-w-[940px]">
          <div className="relative h-[256px] w-[340px]">
            <div className="h-[256px] w-[340px]  bg-gradient-to-br from-primary/10 to-primary/5 rounded-lg ring-1 ring-primary/20 absolute bottom-6 right-6">
            </div>
            <Image src={"/modern.png"} fill={true} alt="info" className="rounded-lg" />
          </div>
          <p className="text-black dark:text-white max-w-[407px] text-[40px] leading-relaxed">Leverage/Monetise your Real World Asset</p>
        </div>
        <p className="text-black dark:text-white leading-relaxed text-center mt-16 w-full text-[30px] max-w-[786px]"><q>
          Globally, many are asset-rich and cash-poor; <span className="text-secondary">PowerDfi</span> is here to balance the playing field.</q></p>
        <h4 className="text-text-grey text-center text-[30px]">-Tim Webb founder CEO-</h4>
      </section>
      <section className="ring-1 ring-primary/40 flex gap-24 justify-center px-8 sm:px-12 md:px-24 py-24">
        <article className="flex-1 max-w-[635px]">
          <h3 className="font-semibold leading-relaxed text-secondary text-[40px]">The PowerDfi platform secures your RWAs and Value on the blockchain</h3>
          <p className="text-[30px] leading-relaxed text-black dark:text-white mt-8">Cyber Security and Asset protection is one of our primary focuses.
            <br />
            PDFI uses state of the art protocols, fraud detection and hacker proofing strategies to keep your digitized asset safe.
            <br />
            PDFI Power Vault is built with the highest encryption and protocols on the market to keep your assets secure on the blockchain.</p>
        </article>
        <div className="relative w-full aspect-square flex-1 h-full justify-end max-w-[536px]">
          <video src="/vid.mp4" controls></video>
        </div>
      </section>
      <section id="services" className="py-12 px-8 sm:px-12 md:px-24 flex flex-col items-center">
        <div className="flex flex-col gap-6 w-full items-start max-w-4xl">
          <div className="text-secondary inline font-semibold p-1.5 rounded-md ring-1 ring-primary/40 w-fit text-[40px]">A plan that is right for you</div>
          <p className="text-black dark:text-white text-[24px] max-w-[970px]">The first three are free if using PDFI leveraging or storing value in our Power vault service, Payment is accepted in stable coins or PDFI token</p>
        </div>
        <div className="flex flex-col gap-16 mt-24 items-center ">
          <div className="flex items-center w-full justify-center max-w-4xl">
            <div className="h-[359px] w-[320px]  flex items-center justify-center  max-w-sm  relative rounded-xl bg-white dark:bg-transparent dark:gradient left-24 z-10 ring-1 ring-primary/20 shadow-sm">
              <div className="relative">
                <Image src={isDarkMode ? "/art1.png" : "/lart1.png"} alt="" height={150} width={150} className="rounded-xl relative" />
              </div>
            </div>
            <article className="w-[820px] h-[456px] flex flex-col justify-center  bg-white dark:gradient ring-1 ring-primary/20 relative rounded-xl p-16 py-0 pl-36">
              <h3 className="text-black dark:text-white text-[30px]">Our tokenisation Plans</h3>
              <h4 className="mt-5 text-[24px] text-black dark:text-white">1st Assets</h4>
              <div className="flex flex-col gap-3 text-black dark:text-white text-[24px] mt-6">
                <div className="flex gap-24">
                  <div>Additional Assets:</div>
                  <div>Free</div>
                </div>
                <div className="flex gap-24">
                  <div>up to $100,000</div>
                  <div>$49 each</div>
                </div>
                <div className="flex gap-24">
                  <div>Above $100,000</div>
                  <div>$99 each</div>
                </div>
              </div>
              <button className="bg-secondary h-[65px] w-[184px] rounded-md mt-10">Try for free</button>
              {/* <p className="text-md text-center text-text-grey self-center mt-8">The first three are free if using PDFI leveraging or storing value in our Power vault service, Payment is accepted in stable coins or PDFI token</p> */}
            </article>
          </div>
          <div className="flex items-center flex-row-reverse max-w-4xl">
            <div className="h-[359px] w-[320px]  flex items-center justify-center  max-w-md aspect-square relative rounded-xl bg-white dark:bg-transparent dark:gradient right-24 z-10 ring-1 ring-primary/20 shadow-sm">
              <div className="relative">
                <Image src={isDarkMode ? "/art2.png" : "/lart2.png"} alt="" height={150} width={150} className="rounded-xl relative" />
              </div>
            </div>
            <article className="w-[820px] h-[456px] bg-white dark:gradient ring-1 ring-primary/20 relative rounded-xl p-16 py-0 pr-36 flex flex-col justify-center">
              <h3 className="text-black dark:text-white text-[30px]">Our Storage Plans</h3>
              <div className="flex flex-col gap-4 text-black dark:text-white text-[24px] mt-6">
                <div>3 months 1 PDFI token per month</div>
                <div>6 months 2 PDFI token per month</div>
                <div>1 year plus 3 PDFI tokens per month</div>
              </div>
              <button className="bg-secondary h-[65px] w-[184px] rounded-md mt-10">Get Started</button>
              {/* <p className="text-md text-center text-text-grey self-center mt-8">The first three are free if using PDFI leveraging or storing value in our Power vault service, Payment is accepted in stable coins or PDFI token</p> */}
            </article>
          </div>
          <div className="flex items-center w-full max-w-4xl">
            <div className="h-[359px] w-[320px]  flex items-center justify-center max-w-xl aspect-square relative rounded-xl bg-white dark:bg-transparent dark:gradient left-24 z-10 ring-1 ring-primary/20 shadow-sm">
              <div className="relative">
                <Image src={isDarkMode ? "/art3.png" : "/lart3.png"} alt="" height={150} width={150} className="rounded-xl relative" />
              </div>
            </div>
            <article className="w-[820px] h-[456px] flex flex-col justify-center bg-white dark:gradient ring-1 ring-primary/20 relative rounded-xl p-16 py-0 pl-36">
              <h3 className="text-black dark:text-white text-[30px]">Leverage/Monetise Plans</h3>
              <h4 className="mt-8 text-[24px] text-black dark:text-white">Swap your RAW token for PDFI  tokens</h4>
              <div className="flex flex-col gap-4 text-black dark:text-white text-[24px] mt-6">
                Offer your RWA token for trade or to collateralize a P2P loan/funding.
              </div>
              <button className="bg-secondary h-[65px] w-[184px] rounded-md mt-10">Get Started</button>
              {/* <p className="text-md text-center text-text-grey self-center mt-8">The first three are free if using PDFI leveraging or storing value in our Power vault service, Payment is accepted in stable coins or PDFI token</p> */}
            </article>
          </div>
        </div>
      </section>
    </main>
  );
}
