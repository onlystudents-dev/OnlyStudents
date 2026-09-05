import {toast} from "react-toastify";

export const ERROR_MESSAGES: Record<string, Record<string, string>> = {
  "en": {
    "WRONG_CREDENTIALS": "The username doesn't pair with the password!",
    "TOO_MANY_REQUESTS": "Too many attempts, try again later.",
    "PASSWORD_WEAK": "Your password is too weak, please try something more secure!",
    "PASSWORD_EXPOSED": "Your password was found in a data breach! Please use another one, and change other accounts which use it.",
    "INVALID_PASSWORD_RESET_TOKEN": "Invalid Password Reset Token!",
    "INVALID_VERIFICATION_CODE": "Invalid Verification Code!"
  },
  "hu": {
    "WRONG_CREDENTIALS": "A belépési adatok helytelenek!",
    "TOO_MANY_REQUESTS": "Túl sok próbálkozás, próbáld újra később.",
    "PASSWORD_WEAK": "A jelszavad túlságosan egyszerű, kérlek próbálj valami biztonságosabbat!",
    "PASSWORD_EXPOSED": "A jelszavadat megtaláltuk egy adatszivárgásban. Kérlek használj egy másikat és változtasd meg máshol is a jelszavaidat, ahol ezt használod!",
    "INVALID_PASSWORD_RESET_TOKEN": "Érvénytelen Jelszó-visszaállítási token!",
    "INVALID_VERIFICATION_CODE": "Érvénytelen kód!"
  },
  "de": {
    "WRONG_CREDENTIALS": "Der Benutzername passt nicht zum Passwort!",
    "TOO_MANY_REQUESTS": "Zu viele Versuche, bitte versuche es später erneut.",
    "PASSWORD_WEAK": "Dein Passwort ist zu schwach, bitte wähle ein sichereres!",
    "PASSWORD_EXPOSED": "Dein Passwort wurde in einem Datenleck gefunden! Bitte verwende ein anderes und ändere es auch überall dort, wo du es benutzt.",
    "INVALID_PASSWORD_RESET_TOKEN": "Ungültiger Token zum Zurücksetzen des Passworts!",
    "INVALID_VERIFICATION_CODE": "Ungültiger Bestätigungscode!"
  },
  "rom": {
    "WRONG_CREDENTIALS": "O anav taj e parola na si khetane!",
    "TOO_MANY_REQUESTS": "But proba sas! Proba pale, kana aresela.",
    "PASSWORD_WEAK": "E parola si but loki! La zoraleder parola.",
    "PASSWORD_EXPOSED": "Arakhlam e parola jekhe databaza! La aver parola taj paruvel la t'e o avera thana, kote koris tut la.",
    "INVALID_PASSWORD_RESET_TOKEN": "Nalačho tokeno e parolakero reseto!",
    "INVALID_VERIFICATION_CODE": "Nalačho kod e verifikacia!"
  },
  "ro": {
    "WRONG_CREDENTIALS": "Numele de utilizator nu se potrivește cu parola!",
    "TOO_MANY_REQUESTS": "Prea multe încercări! Încercă din nou mai târziu.",
    "PASSWORD_WEAK": "Parola ta este prea slabă! Te rog, alege una parola mai sigură.",
    "PASSWORD_EXPOSED": "Parola ta a fost găsită într-o scurgere de date! Te rog, folosește alta și schimb-o și în alte conturi unde o folosești.",
    "INVALID_PASSWORD_RESET_TOKEN": "Token invalid pentru resatarea parolei!",
    "INVALID_VERIFICATION_CODE": "Cod de verificar invalid!"
  },
  "sk": {
    "WRONG_CREDENTIALS": "Používateľské meno sa nezhoduje s heslom!",
    "TOO_MANY_REQUESTS": "Príliš veľa pokusov! Skúste znova neskôr.",
    "PASSWORD_WEAK": "Vaše heslo je príliš slabé! Skúste bezpečnejšie heslo.",
    "PASSWORD_EXPOSED": "Vaše heslo sa našlo v úniku údajov! Použite iné heslo a zmeňte ho aj v iných kontách, kde ho používate.",
    "INVALID_PASSWORD_RESET_TOKEN": "Neplatný token na zresetovanie hesla!",
    "INVALID_VERIFICATION_CODE": "Neplatný verifikačný kód!"
  },
  "sr": {
    "WRONG_CREDENTIALS": "Korisničko ime se ne poklapa sa lozinkom!",
    "TOO_MANY_REQUESTS": "Previše pokušaja! Pokušajte ponovo kasnije.",
    "PASSWORD_WEAK": "Vaša lozinka je preslaba! Pokušajte nešto sigurnije.",
    "PASSWORD_EXPOSED": "Vaša lozinka je pronađena u curenju podataka! Koristite drugu i promenite je svuda gde je koristite.",
    "INVALID_PASSWORD_RESET_TOKEN": "Nevažeći token za resetovanje lozinke!",
    "INVALID_VERIFICATION_CODE": "Nevažeći verifikacioni kod!"
  },
  "uk": {
    "WRONG_CREDENTIALS": "Ім'я користувача не відповідає паролю!",
    "TOO_MANY_REQUESTS": "Занадто багато спроб! Спробуйте пізніше.",
    "PASSWORD_WEAK": "Ваш пароль занадто слабкий! Оберіть надійніший.",
    "PASSWORD_EXPOSED": "Ваш пароль знайшли в витіку даних! Оберіть інший і замініть его і в інших сервисах, де ви его використовуєте.",
    "INVALID_PASSWORD_RESET_TOKEN": "Недійсний токен для відновлення пароля!",
    "INVALID_VERIFICATION_CODE": "Недійсний код підтвердження!"
  },
  "zh": {
    "WRONG_CREDENTIALS": "用户名与密码不匹配！",
    "TOO_MANY_REQUESTS": "尝试次数过多，请稍后再试。",
    "PASSWORD_WEAK": "您的密码太弱，请尝试更安全的密码！",
    "PASSWORD_EXPOSED": "您的密码已在数据泄露中被发现！请更换密码，并同时修改其他使用该密码的账户。",
    "INVALID_PASSWORD_RESET_TOKEN": "无效的密码重置令牌！",
    "INVALID_VERIFICATION_CODE": "无效的验证码！"
  },
  "zh-TW": {
    "WRONG_CREDENTIALS": "使用者名稱與密碼不匹配！",
    "TOO_MANY_REQUESTS": "嘗試次數過多，請稍後再試。",
    "PASSWORD_WEAK": "您的密碼太弱，請嘗試更安全的密碼！",
    "PASSWORD_EXPOSED": "您的密碼已在資料外洩中被發現！請更換密碼，並同時修改其他使用該密碼的帳戶。",
    "INVALID_PASSWORD_RESET_TOKEN": "無效的密碼重設令牌！",
    "INVALID_VERIFICATION_CODE": "無效的驗證碼！"
  }
};

const KNOWN_ERRORS = new Set(Object.keys(ERROR_MESSAGES.en));

export function showError(input: Record<string, string>) {
  const code = input.error;

  if (
    Object.keys(input).length !== 1 ||
    typeof code !== "string" ||
    !KNOWN_ERRORS.has(code)
  ) {
    console.error("[showError] rejected:", input);
    return;
  }

  const locale = navigator.language ?? "en";
  toast.error(ERROR_MESSAGES[locale]?.[code] ?? ERROR_MESSAGES.en[code]);
}
