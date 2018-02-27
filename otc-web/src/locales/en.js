export default {
  header: {
    navigation: {
      whitepapers: 'Whitepapers',
      downloads: 'Downloads',
      explorer: 'Explorer',
      blog: 'Blog',
      roadmap: 'Roadmap',
    },
  },
  footer: {
    getStarted: 'Get started',
    explore: 'Explore',
    community: 'Community',
    wallet: 'Get Wallet',
    infographics: 'Infographics',
    whitepapers: 'Whitepapers',
    blockchain: 'Blockchain Explorer',
    blog: 'Blog',
    twitter: 'Twitter',
    reddit: 'Reddit',
    github: 'Github',
    telegram: 'Telegram',
    slack: 'Slack',
    roadmap: 'Roadmap',
    skyMessenger: 'Sky-Messenger',
    cxPlayground: 'CX Playground',
    team: 'Team',
    subscribe: 'Mailing List',
    market: 'Markets',
    bitcoinTalks: 'Bitcointalks ANN',
    instagram: 'Instagram',
    facebook: 'Facebook',
    discord: 'Discord',
  },
  distribution: {
    rate: 'Current OTC rate: {rate} SKY/BTC',
    inventory: 'Current inventory: {coins} SKY available',
    title: 'Skycoin OTC',
    heading: 'Skycoin OTC',
    headingEnded: 'Skycoin OTC is currently closed',
    ended: `<p>Join the <a href="https://t.me/skycoin">Skycoin Telegram</a>
      or follow the
      <a href="https://twitter.com/skycoinproject">Skycoin Twitter</a>.`,
    instructions: `<p>You can check the current market value for <a href="https://coinmarketcap.com/currencies/skycoin/">Skycoin at CoinMarketCap</a>.</p>

<p>To use the Skycoin OTC:</p>

<ul>
  <li>Enter your Skycoin address below</li>
  <li>You&apos;ll receive a unique Bitcoin address to purchase SKY</li>
  <li>Send BTC to the address</li>
</ul>

<p>You can check the status of your order by entering your address and selecting <strong>Check status</strong>.</p>
<p>Each time you select <strong>Get Address</strong>, a new BTC address is generated. A single SKY address can have up to 5 BTC addresses assigned to it.</p>
    `,
    statusFor: 'Status for {skyAddress}',
    enterAddress: 'Enter Skycoin address',
    getAddress: 'Get address',
    checkStatus: 'Check status',
    loading: 'Loading...',
    btcAddress: 'BTC address',
    errors: {
      noSkyAddress: 'Please enter your SKY address.',
      coinsSoldOut: 'Skycoin OTC is currently sold out, check back later.',
    },
    statuses: {
      waiting_deposit: '[tx-{id} {updated}] Waiting for BTC deposit.',
      waiting_send: '[tx-{id} {updated}] BTC deposit confirmed. Skycoin transaction is queued.',
      waiting_confirm: '[tx-{id} {updated}] Skycoin transaction sent.  Waiting to confirm.',
      done: '[tx-{id} {updated}] Completed. Check your Skycoin wallet.',
    },
  },
};
